//Package main
//
//walk_single однопоточный поиск дублей
//поиск осуществляется по наименованию и размеру файла
//walk_multi многопоточный поиск дублей
//поиск осуществляется по наименованию и размеру файла
//main Основная функция для запуска поиска дублей
//Флаги при запуске :
// --delete - удаление дублей с оставлением одного экземпляра файла
// --multi - запуск многопоточного поиска
// --path -  путь до верхней папки ,в которой будет осуществлен поиск дублей
package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	//"log"
	"sort"
	"sync"
	"time"
)

var walked_files = make(map[string][]string)
var walked_lock sync.Mutex

// walk_single однопоточный поиск дублей
//поиск осуществляется по наименованию и размеру файла
func walk_single(path string) {
	log.SetFormatter(&log.JSONFormatter{})
	//hlog := log.WithFields(standardFields)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "walk_single",
			"error": err,
			"path":  path,
		}).Error("not find path")
		//log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if err != nil {
				log.WithFields(log.Fields{
					"func":  "walk_single",
					"error": err,
					"dir":   path,
				}).Panic("error dir")
			}
			walk_single(path + "\\" + f.Name())
		} else {
			hash := fmt.Sprintf("%s%d", f.Name(), f.Size())
			log.WithFields(log.Fields{
				"func":     "walk_single",
				"file":     f,
				"nameFile": f.Name(),
				"sizeEile": f.Size(),
			}).Info("info files")
			walked_files[hash] = append(walked_files[hash], path+"\\"+f.Name())

		}

	}
}

//walk_multi многопоточный поиск дублей
//поиск осуществляется по наименованию и размеру файла
func walk_multi(wg *sync.WaitGroup, path string) {
	log.SetFormatter(&log.JSONFormatter{})
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "walk_multi",
			"error": err,
			"path":  path,
		}).Error("not find path")
		//	log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			wg.Add(1)
			if err != nil {
				log.WithFields(log.Fields{
					"func":  "walk_multi",
					"error": err,
					"dir":   path,
				}).Panic("error dir")
			}
			go walk_multi(wg, path+"\\"+f.Name())
		} else {
			hash := fmt.Sprintf("%s%d", f.Name(), f.Size())
			log.WithFields(log.Fields{
				"func":     "walk_multi",
				"file":     f,
				"nameFile": f.Name(),
				"sizeEile": f.Size(),
			}).Info("info files")
			walked_lock.Lock()
			walked_files[hash] = append(walked_files[hash], path+"\\"+f.Name())
			walked_lock.Unlock()
		}
	}
	wg.Done()
}

//main Основная функция для запуска поиска дублей
//Флаги при запуске :
// --delete - удаление дублей с оставлением одного экземпляра файла
// --multi - запуск многопоточного поиска
// --path -  путь до верхней папки ,в которой будет осуществлен поиск дублей
func main() {
	var arg_delete bool
	var arg_multi bool
	var arg_help bool
	var arg_path string
	flag.BoolVar(&arg_delete, "delete", false, "delete duplicate files")
	flag.BoolVar(&arg_multi, "multi", false, "run program in miltithreaded mode")
	flag.BoolVar(&arg_help, "help", false, "about programm")
	flag.BoolVar(&arg_help, "h", false, "about programm")
	flag.StringVar(&arg_path, "path", ".", "path to found duplicate files")
	flag.Parse()
	if arg_help {
		fmt.Println("Программа поиска дублей в директории рекурсивно")
		return
	}
	time1 := time.Now()
	if arg_multi {
		var wg sync.WaitGroup
		wg.Add(1)
		walk_multi(&wg, arg_path)
		wg.Wait()
	} else {
		walk_single(arg_path)
	}
	time2 := time.Now()

	hashes := make([]string, 0, len(walked_files))
	for hash := range walked_files {
		hashes = append(hashes, hash)
	}
	sort.Strings(hashes)

	for _, hash := range hashes {
		if len(walked_files[hash]) > 1 {
			fmt.Printf("%d duplicates:\n", len(walked_files[hash]))
			for _, fn := range walked_files[hash] {
				fmt.Printf(" - %s\n", fn)
			}
			if arg_delete {
				for i, fn := range walked_files[hash] {
					if i == 1 {
						fmt.Printf(" - main files %s\n", fn)
					} else {
						os.Remove(fn)
						fmt.Printf(" - delete files %s\n", fn)
					}

				}
			}
		}
	}
	fmt.Printf("done in %s\n", time2.Sub(time1).String())
}
