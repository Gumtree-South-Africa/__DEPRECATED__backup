package work

//import (
//
//)

func (Mongo Mongo) dump(host string, port string, username string, tableName string, destination string, errc chan error) Input {
	var input Input
	return input
}

//func (Mongo Mongo) rsync(source string, destination string, errc chan error) {
//	if len(source) > 0 && len(destination) > 0 {
//		cmd := exec.Command("/usr/bin/rsync", "-avx", "-e", "\"ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null\"", source, destination)
//		out, err := os.Create("/tmp/rsync.log")
//		cmd.Stdout = out
//		cmd.Stderr = os.Stderr
//		if err != nil {
//			log.Fatal("Error while creating error log", err)
//			errc <- err
//		}
//		execErr := cmd.Run()
//		if execErr != nil {
//			log.Fatal("Error while performing rsync", execErr)
//			errc <- execErr
//		}
//		log.Println(source, destination)
//	}
//}