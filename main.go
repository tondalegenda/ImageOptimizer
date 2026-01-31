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

var fileSizeLimit int = 1000
var sideTooBigLimit int = 1500
var shorterSideTarget int = 600

func main() {
	files, err := os.ReadDir(".") //vytvoří proměnné files a err, které koukají do složky ve které jsme
	if err != nil {
		fmt.Println("Něco je špatně?")
	}

	images := findImagesInFolder(files)
	printFullImageTable(images)
	//askToOptimize()

	for {
		MenuChoice := showMainMenu()
		if MenuChoice == 0 {
			fmt.Println("Sbohem")
			break
		}
		switch MenuChoice {
		case 1: // automatický resize všech problematických
			targets := selectOptiImagesAll(images)
			optimizeSelectedImages(targets)
		case 2: //resize jen těch co mají moc kb
			targets := selectOptiImagesFileSize(images)
			optimizeSelectedImages(targets)
		case 3: //resize jen těch co mají velké rozměry
			targets := selectOptiImagesDimensions(images)
			optimizeSelectedImages(targets)
		case 4: //konverze PNG obrázků na JPG
			targets := selectConvImagesPng(images)
			convertSelectedImages(targets)
		case 5: //konverze JPG obrázků na PNG
			targets := selectConvImagesJpeg(images)
			convertSelectedImages(targets)
		case 99:
			showChangeLimitMenu()
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
			isTooBig := kbSize > int64(fileSizeLimit)
			isTooWide := config.Width > sideTooBigLimit || config.Height > sideTooBigLimit
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

func showMainMenu() int {
	fmt.Println("===== Chceš optimalizovat obrázky? =====")
	fmt.Println("(1) Optimalizovat vše")
	fmt.Println("(2) Optimalizovat jen velikost")
	fmt.Println("(3) Optimalizovat jen rozměry")
	fmt.Println("(4) Konvertovat obrázky na .jpg")
	fmt.Println("(5) Konvertovat obrázky na .png")
	fmt.Println("(99) Změnit limity")
	fmt.Println("(0) Zpět")
	var menuchoice int
	fmt.Scan(&menuchoice)
	return menuchoice
}

func showChangeLimitMenu() int {
	fmt.Println("===== Nastavené limity =====")
	fmt.Println("Maximální souborová velikost:", fileSizeLimit, "kb \t\t (1) Změnit")
	fmt.Println("Délka strany moc velkých obrázků:", sideTooBigLimit, "px \t (2) Změnit")
	fmt.Println("Cílová velikost kratší strany:", shorterSideTarget, "px \t\t (3) Změnit")
	fmt.Println("(0) Zpět")
	var menuchoice int
	var limitInput int
	fmt.Scan(&menuchoice)
	switch menuchoice {
	case 1:
		fmt.Println("Zadej novou maximální souborovou velikost")
		fmt.Scan(&limitInput)
		if limitInput > 0 && limitInput < 10000 {
			fileSizeLimit = limitInput
			fmt.Println("Maximální souborová velikost změněna na", fileSizeLimit, "kb")
		} else {
			fmt.Println("Neplatná hodnota")
		}
	case 2:
		fmt.Println("Zadej novou minimální délku strany pro vyhodnocení velkých obrázků")
		if limitInput > 0 && limitInput < 20000 {
			fmt.Scan(&limitInput)
			sideTooBigLimit = limitInput
			fmt.Println("Délka strany moc velkých obrázků změněna na", sideTooBigLimit, "px")
		} else {
			fmt.Println("Neplatná hodnota")
		}
	case 3:
		fmt.Println("Zadej novou cílovou velikost kratší strany")
		if limitInput > 0 && limitInput < 2000 {
			fmt.Scan(&limitInput)
			shorterSideTarget = limitInput
			fmt.Println("Cílová velikost kratší strany změněna na", shorterSideTarget, "px")
		} else {
			fmt.Println("Neplatná hodnota")
		}
	default:
		fmt.Println("Neplatná hodnota v menu")
	}
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
	for _, OptiImagesFileSizeList := range ImageList {
		if OptiImagesFileSizeList.FileSizeError {
			SelectedImages = append(SelectedImages, OptiImagesFileSizeList)
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func selectOptiImagesDimensions(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, OptiImagesDimensionsList := range ImageList {
		if OptiImagesDimensionsList.DimensionsError {
			SelectedImages = append(SelectedImages, OptiImagesDimensionsList)
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func selectConvImagesJpeg(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, ConvImagesJpegList := range ImageList {
		extension := filepath.Ext(ConvImagesJpegList.Name)
		if extension == ".jpg" || extension == ".JPG" {
			SelectedImages = append(SelectedImages, ConvImagesJpegList)
		} else {
			continue
		}
	}
	printFullImageTable(SelectedImages)
	return SelectedImages
}

func selectConvImagesPng(ImageList []ImageFile) []ImageFile {
	var SelectedImages []ImageFile
	for _, ConvImagesPngList := range ImageList {
		extension := filepath.Ext(ConvImagesPngList.Name)
		if extension == ".png" || extension == ".PNG" {
			SelectedImages = append(SelectedImages, ConvImagesPngList)
		} else {
			continue
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
			newImage = resize.Resize(0, uint(shorterSideTarget), decodedImage, resize.Lanczos3)
			fmt.Println("Snižuji výšku na 600px")
			openAndSaveImages(newImage, selectedImage)
		} else if selectedImage.Width < selectedImage.Height && selectedImage.Width > 600 {
			newImage = resize.Resize(uint(shorterSideTarget), 0, decodedImage, resize.Lanczos3)
			fmt.Println("Snižuji šířku na 600px")
			openAndSaveImages(newImage, selectedImage)
		} else {
			fmt.Println("S timhle nevim co dělat tvl")
		}

	}
}

func convertSelectedImages(ImageList []ImageFile) {
	for _, selectedImage := range ImageList {
		fmt.Println("Konvertuji", selectedImage.Name)
		file, _ := os.Open(selectedImage.Name)   //otevře file
		decodedImage, _, _ := image.Decode(file) //uloží file do paměti
		file.Close()                             //zavře file

		convertImageFormat(decodedImage, selectedImage)

	}
}

func convertImageFormat(newImage image.Image, selectedImage ImageFile) {
	extension := filepath.Ext(selectedImage.Name) // zjistí typ obrázku a encodne
	if extension == ".jpg" || extension == ".JPG" {
		newFileName := selectedImage.Name[:len(selectedImage.Name)-len(extension)] + ".png"
		newFile, err := os.Create("conv_" + newFileName)
		fmt.Println(selectedImage.Name, "se bude nově jmenovat", newFileName)
		selectedImage.Name = newFileName
		if err != nil {
			fmt.Println("Přeskakuji")
			return
		}
		png.Encode(newFile, newImage)
	} else {
		newFileName := selectedImage.Name[:len(selectedImage.Name)-len(extension)] + ".jpg"
		newFile, err := os.Create("conv_" + newFileName)
		fmt.Println(selectedImage.Name, "se bude nově jmenovat", newFileName)
		selectedImage.Name = newFileName
		if err != nil {
			fmt.Println("Přeskakuji")
			return
		}
		jpeg.Encode(newFile, newImage, nil)
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
