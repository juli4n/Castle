// A Go mirror of libfuse's hello.c

package main

import (
	"castle/backup"
	"castle/blob"
	"castle/object"
	"flag"
	"github.com/hanwen/go-fuse/fuse"
	"log"
	"strings"
)

type BackupFs struct {
	fuse.DefaultFileSystem
	rootNode      *backup.Tree
	objectStorage object.ObjectStorage
	blobStorage   blob.BlobStorage
}

func (self *BackupFs) init(name string) {
	var root []byte
	var err error
	if root, err = self.blobStorage.Get(name); err != nil {
		log.Fatalf("Error loading %s information.\n", name)
	}
	lastBackup := &backup.Backup{}
	self.objectStorage.Get(string(root), lastBackup)

	log.Printf("Mounting backup: %s", lastBackup.Name)
	self.rootNode = &backup.Tree{}
	self.objectStorage.Get(lastBackup.Current, self.rootNode)
	log.Printf("Current Node: %v", self.rootNode)
}

func (self *BackupFs) getDescendantTree(name string) (child *backup.TreeNode) {
	dirTree := *self.rootNode
	dirs := strings.Split(name, "/")
	for i := range dirs {
		if child = dirTree.FindChildNode(dirs[i]); child != nil {
			self.objectStorage.Get(child.ObjRef, &dirTree)
		} else {
			return nil
		}
	}
	return child
}

func (self *BackupFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	if name == "" {
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755}, fuse.OK
	}
	dirNode := self.getDescendantTree(name)
	if dirNode == nil {
		return nil, fuse.ENOENT
	}
	if dirNode.IsDir {
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0644}, fuse.OK
	}
	data, _ := self.blobStorage.Get(dirNode.ObjRef)
	size := uint64(len(data))
	return &fuse.Attr{Mode: fuse.S_IFREG | 0755, Size: size}, fuse.OK
}

func (self *BackupFs) buildDirChannel(dir *backup.Tree) chan fuse.DirEntry {
	c := make(chan fuse.DirEntry, len(dir.Childs))
	for i := range dir.Childs {
		childNode := dir.Childs[i]
		if childNode.IsDir {
			c <- fuse.DirEntry{Name: childNode.Name, Mode: fuse.S_IFDIR}
		} else {
			c <- fuse.DirEntry{Name: childNode.Name, Mode: fuse.S_IFREG}
		}
	}
	close(c)
	return c
}

func (self *BackupFs) OpenDir(name string, context *fuse.Context) (c chan fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c := self.buildDirChannel(self.rootNode)
		return c, fuse.OK
	} else {
		dirNode := self.getDescendantTree(name)
		if dirNode == nil {
			return nil, fuse.ENOENT
		}
		dirTree := &backup.Tree{}
		self.objectStorage.Get(dirNode.ObjRef, dirTree)
		return self.buildDirChannel(dirTree), fuse.OK
	}
	return nil, fuse.ENOENT
}

func (self *BackupFs) Open(name string, flags uint32, context *fuse.Context) (file fuse.File, code fuse.Status) {
	child := self.getDescendantTree(name)
	if child == nil {
		return nil, fuse.ENOENT
	}
	if data, err := self.blobStorage.Get(child.ObjRef); err != nil {
		log.Fatalf("Internal error: ObjRef %s is a dangling pointer", child.ObjRef)
	} else {
		return fuse.NewDataFile(data), fuse.OK
	}
	return nil, fuse.ENOENT
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("Usage:\n  fuse mountpoint backup-folder")
	}
	var blobStorage blob.BlobStorage
	var objectStorage object.ObjectStorage
	var err error
	storageFile := flag.Arg(1) + "/.objects"
	if blobStorage, err = blob.NewFileBasedBlobStorage(storageFile); err != nil {
		log.Fatalf("Error loading object database. Are you running from the correct directory? (%s)\n", storageFile)
	}
	objectStorage = object.NewJSONStorage(blobStorage)
	backupFs := &BackupFs{blobStorage: blobStorage, objectStorage: objectStorage}
	backupFs.init("HEAD")
	nfs := fuse.NewPathNodeFs(backupFs, nil)
	state, _, err := fuse.MountNodeFileSystem(flag.Arg(0), nfs, nil)
	if err != nil {
		log.Fatal("Mount fail: %v\n", err)
	}
	state.Loop()
}
