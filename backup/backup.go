package backup

import (
  "castle/object"
  "castle/blob"
  "io/ioutil"
  "os"
  "fmt"
  "time"
  "container/list"
)

// A tree of references to persisted objects.
type Tree struct {
	// Child nodes.
	Childs []TreeNode
}

// 
type TreeNode struct {
	// The name of the child.
	Name string
	// The id of the referenced object.
	ObjRef string
	// true iff is a leaf.
	IsDir bool
}

type Backup struct {
  Name string
  Date int64
  Previous string
  Current string
}

func CreateBackup(name string, previous string, rootDir string, blobStorage blob.BlobStorage, gobStorage object.ObjectStorage) (backupRef string, err error) {
  var treeRef string
  if treeRef, err = Walk(rootDir, blobStorage, gobStorage); err != nil {
    return "", err
  }
  backup := &Backup{name, time.Now().Unix(), previous, treeRef}
  if backupRef, err = gobStorage.Put(backup); err != nil {
    return "", err
  }
  return backupRef, nil
}

func GetBackupHistory(backupRef string, gobStorage object.ObjectStorage) (history []Backup, err error) {
  historyList := list.New()
  var backup *Backup
  for prev := backupRef; backup == nil || backup.Previous != ""; prev = backup.Previous {
    backup = &Backup{}
    if err = gobStorage.Get(prev, backup); err != nil {
      return nil, err
    }
    historyList.PushBack(backup)
  }
  history = make([]Backup, historyList.Len())
  i := 0
  for v := historyList.Front(); v != nil;  v = v.Next() {
    history[i] = *(v.Value.(*Backup))
    i++
  }
  return history, nil
}

func Walk(rootDir string, blobStore blob.BlobStorage, gobStore object.ObjectStorage) (objRef string, err error) {

  var files []os.FileInfo

  if files, err = ioutil.ReadDir(rootDir); err != nil {
    return "", fmt.Errorf("Error reading content of directory %s", rootDir)
  }
  t := &Tree{make([]TreeNode, len(files))}
  for i := range(files) {
    f := files[i]
    filename := rootDir + "/" + f.Name()
    if !f.IsDir() {
      var fileContent []byte
      if fileContent, err = ioutil.ReadFile(filename); err != nil {
        return "", fmt.Errorf("Error reading content of file %s", filename)
      }
      if objRef, err = blobStore.Put(fileContent); err != nil {
        return "", fmt.Errorf("Error saving blob object for file %s", filename)
      }
      t.Childs[i] = TreeNode{f.Name(), objRef, false}
    } else {
      if objRef, err = Walk(filename, blobStore, gobStore); err != nil {
        return "", err
      }
      t.Childs[i] = TreeNode{f.Name(), objRef, true}
    }
  }
  if objRef, err = gobStore.Put(t); err != nil {
    return "", err
  }
  return objRef, nil
}

func ReconstructTree(rootObjRef string, outputFolder string, blobStore blob.BlobStorage, gobStore object.ObjectStorage) (err error) {
  var t *Tree = &Tree{}
  if err = gobStore.Get(rootObjRef, t); err != nil {
    return err//errors.New("Reference %s not found. Check your configuration. The repository seems to be broken.")
  }
  os.Mkdir(outputFolder, os.ModeDir | 0700)
  for i := range(t.Childs) {
    child := t.Childs[i]
    if child.IsDir {
      dirName := outputFolder + "/" + child.Name
      ReconstructTree(child.ObjRef, dirName, blobStore, gobStore)
    } else {
      var file *os.File
      var content []byte
      filename := outputFolder + "/" + child.Name
      if file, err = os.Create(filename); err != nil {
        return err
      }
      defer file.Close()
      if content, err = blobStore.Get(child.ObjRef); err != nil {
        return err
      }
      if _, err = file.Write(content); err != nil {
        return err
      }
    }
  }
  return nil
}



