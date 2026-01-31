package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type ImageFile struct {
	Name            string
	Size            int64
	Width           int
	Height          int
	FileSizeError   bool
	DimensionsError bool
}

func main() {
	files, err := os.ReadDir(".") //vytvoří proměnné files a err, které koukají do složky ve které jsme
	if err != nil {
		fmt.Println("Něco je špatně?")
	}

	images := findImagesInFolder(files)
	printFullImageTable(images)
	//askToOptimize()

	for {
		MenuChoice := askToOptimize()
		if MenuChoice == 0 {
			fmt.Println("Sbohem")
			break
		}
		switch MenuChoice {
		case 1:
			targets := selectOptiImagesAll(images)
			optimizeSelectedImages(targets)
		case 2:
			targets := selectOptiImagesFileSize(images)
			optimizeSelectedImages(targets)
		case 3:
			targets := selectOptiImagesDimensions(images)
			optimizeSelectedImages(targets)
		default:
			fmt.Println("Neplatná volba, zkus to znovu.")
		}
	}

}

func findImagesInFolder(filesInDirectory []os.DirEntry) []ImageFile {
	var images []ImageFile //vytvoří proměnnou images, která je pole (slice) struktur, kterou jsem definoval výš

	for _, file := range filesInDirectory { //Pro každý prvek (který si pojmenujeme file) ze seznamu files udělej následující...
		//fmt.Println(file.Name()) 					//vypiš jméno souboru - tohle už zakomentujeme, neobrázkové soubory nezájem
		if filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".JPG" || filepath.Ext(file.Name()) == ".PNG" { //pokud je přípona souboru jpg nebo png? tohle musí jít zapsat líp ale zatím jsem prý na to blbej
			info, _ := file.Info()                //prommná info do které dám info o obrázku. _ je tam, protože file.Info() vrací info A ERROR jako všechno v go, takže nemusím řešit errory
			f, _ := os.Open(file.Name())          //proměnná f ve které otevřu obrázek
			config, _, _ := image.DecodeConfig(f) //proměnná config, která přečte metadata z hlavičky
			f.Close()                             //zavři obrázek
			kbSize := info.Size() / 1024
			isTooBig := kbSize > 1000
			isTooWide := config.Width > 1500 || config.Height > 1500
			newImage := ImageFile{Name: file.Name(), Size: kbSize, Width: config.Width, Height: config.Height, FileSizeError: isTooBig, DimensionsError: isTooWide} //proměnná newImage, která do struktury definované výš přidá pro každý obrázek údaje o něm
			images = append(images, newImage)                                                                                                                       //přidá strukturu do pole images
		}
	}

	fmt.Println("Našel jsem", len(images), "obrázků")
	fmt.Println(images)
	return images
}

func printFullImageTable(ImageList []ImageFile) {
	fmt.Println("========== SEZNAM SOUBORŮ ==========")
	fmt.Println("Název \t Velikost \t Rozměry \t Optimalizace")
	for _, img := range ImageList {
		var OptimizationReason string
		if img.FileSizeError == true && img.DimensionsError == true {
			OptimizationReason = "Velikost, Rozměry"
		} else if img.FileSizeError == true && img.DimensionsError == false {
			OptimizationReason = "Velikost"
		} else if img.FileSizeError == false && img.DimensionsError == true {
			OptimizationReason = "Rozměry"
		} else {
			OptimizationReason = "OK"
		}
		fmt.Println(img.Name, "\t", img.Size, "\t", img.Width, "x", img.Height, "\t", OptimizationReason)
	}
}

func askToOptimize() int {
	fmt.Println("Chceš optimalizovat obrázky?")
	fmt.Println("(1) Optimalizovat vše")
	fmt.Println("(2) Optimalizovat jen velikost")
	fmt.Println("(3) Optimalizovat jen rozměry")
	fmt.Println("(0) Zpět")
	var menuchoice int
	fmt.Scan(&menuchoice)
	return menuchoice
}

func selectOptiImagesAll(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, OptiImagesAllList := range ImageList {
		if OptiImagesAllList.DimensionsError || OptiImagesAllList.FileSizeError {
			SelectedImages = append(SelectedImages, OptiImagesAllList)
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func selectOptiImagesFileSize(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, OptiImagesAllList := range ImageList {
		if OptiImagesAllList.FileSizeError {
			SelectedImages = append(SelectedImages, OptiImagesAllList)
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func selectOptiImagesDimensions(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, OptiImagesAllList := range ImageList {
		if OptiImagesAllList.DimensionsError {
			SelectedImages = append(SelectedImages, OptiImagesAllList)
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func optimizeSelectedImages(ImageList []ImageFile) {
	for _, selectedImage := range ImageList {
		fmt.Println("Optimalizuji", selectedImage.Name)
		file, _ := os.Open(selectedImage.Name)   //otevře file
		decodedImage, _, _ := image.Decode(file) //uloží file do paměti
		file.Close()                             //zavře file

		var newImage image.Image
		if selectedImage.Width >= selectedImage.Height && selectedImage.Height > 600 {
			newImage = resize.Resize(0, 600, decodedImage, resize.Lanczos3)
			fmt.Println("Snižuji výšku na 600px")
			openAndSaveImages(newImage, selectedImage)
		} else if selectedImage.Width < selectedImage.Height && selectedImage.Width > 600 {
			newImage = resize.Resize(600, 0, decodedImage, resize.Lanczos3)
			fmt.Println("Snižuji šířku na 600px")
			openAndSaveImages(newImage, selectedImage)
		} else {
			fmt.Println("S timhle nevim co dělat tvl")
		}

	}
}

func openAndSaveImages(newImage image.Image, selectedImage ImageFile) {
	newFile, err := os.Create("resized_" + selectedImage.Name) //vytvoří nový soubor s příponou resized_
	if err != nil {
		fmt.Println("Přeskakuji")
		return
	}

	defer newFile.Close()

	extension := filepath.Ext(selectedImage.Name) // zjistí typ obrázku a encodne
	if extension == ".jpg" || extension == ".JPG" {
		jpeg.Encode(newFile, newImage, nil)
	} else {
		png.Encode(newFile, newImage)
	}
}
