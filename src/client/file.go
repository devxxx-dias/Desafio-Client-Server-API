package cliente

import (
	"fmt"
	"os"
)

type File struct {
	Name      string
	Extension string
	MainFile  *os.File
}

func NewFile(name, extension string) (File, error) {
	return File{
		Name:      name,
		Extension: extension,
	}, nil
}

func (f *File) CreateFile() {
	file, err := os.Create(fmt.Sprintf("%s%s", f.Name, f.Extension))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating file: %v\n", err)
	}
	defer file.Close()
}

//func (f *File) WriteFile(res AwesomeApiRequest) {
//	_, err := f.MainFile.WriteString(fmt.Sprintf("%s - %s", res.USDBRL, res.USDBRL.Bid))
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Error while writing file: %v\n", err)
//	}
//	fmt.Println("Quotation saved on file successfully!")
//	fmt.Println("Quotation value: ", res.USDBRL.Bid)
//}
