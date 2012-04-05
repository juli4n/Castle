package main

import (
  "castle/backup"
)

func main() {
/*  var blobStorage blob.BlobStorage
  var objectStorage object.ObjectStorage
  var err error
  switch(os.Args[1]) {
  case "init":
    fmt.Print("Creating a new object database...\n")
    blobStorage, _ := blob.NewFileBasedBlobStorage(".objects")
    blobStorage.PutWithId("source", []byte(os.Args[2]))
    blobStorage.PutWithId("HEAD", []byte(""))
    fmt.Printf("Saving source directory: %s\n", os.Args[2])
  case "backup":
    var backupDir, previous []byte
    var objRef string
    if blobStorage, err = blob.NewFileBasedBlobStorage(".objects"); err != nil {
      fmt.Print("Error loading object database. Are you running %s from the correct directory?\n")
      os.Exit(-1)
    }
    objectStorage = object.NewJSONStorage(blobStorage)
    if backupDir, err = blobStorage.Get("source"); err != nil {
      fmt.Printf("Error loading source information. %s\n", err.Error())
      os.Exit(-1)
    }
    if previous, err = blobStorage.Get("HEAD"); err != nil {
      fmt.Print("Error loading HEAD information.\n")
      os.Exit(-1)
    }
    fmt.Printf("Starting backup of %s\n", string(backupDir))
    if objRef, err = backup.CreateBackup(os.Args[2], string(previous), string(backupDir), blobStorage, objectStorage); err != nil {
      fmt.Printf("Error creating backup: %s.\n", err.Error())
      os.Exit(-1)
    }
    blobStorage.PutWithId("HEAD", []byte(objRef))
    fmt.Printf("Backup completed. Backup %s\n", objRef)
  case "log":
    if blobStorage, err = blob.NewFileBasedBlobStorage(".objects"); err != nil {
      fmt.Print("Error loading object database. Are you running %s from the correct directory?\n")
      os.Exit(-1)
    }
    objectStorage = object.NewJSONStorage(blobStorage)
    var head []byte
    var history []backup.Backup
    if head, err = blobStorage.Get("HEAD"); err != nil {
      fmt.Print("Error loading HEAD information.\n")
      os.Exit(-1)
    }
    if history, err = backup.GetBackupHistory(string(head), objectStorage); err != nil {
      fmt.Printf("Error loading backup log. %s\n", err.Error())
      os.Exit(-1)
    }
    for _, v := range(history) {
      fmt.Printf("Name: %s\nId: %s\nDate: %s\n\n", v.Name, v.Current, time.Unix(v.Date,0))
    }
  case "migrate":
    if blobStorage, err = blob.NewFileBasedBlobStorage(".objects"); err != nil {
      fmt.Print("Error loading object database. Are you running %s from the correct directory?\n")
      os.Exit(-1)
    }
    objectStorage = object.NewJSONStorage(blobStorage)
    if err = backup.ReconstructTree(os.Args[2], os.Args[3], blobStorage, objectStorage); err != nil {
      fmt.Printf("Error migrating to version %s", os.Args[2])
      os.Exit(-1)
    }
    fmt.Print("Migration completed\n")
  }
  */
  backup.Main()
}
