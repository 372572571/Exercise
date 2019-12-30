package diroperation

import (
	"A.brikc/app/toolboxegexpfilter"
	"os"
	"errors"
)
// dir new delete rename

// new dir
// fail return error . succeed return nil
func Newdir(dirname,dirpath string)error{
	if !regexpfilter.MatchAlnum(dirname){
		return errors.New("not [:alnum:].")
	}
	err:=os.MkdirAll(dirpath+dirname,0777)
	if err!=nil {
		return err
	}
	return nil
}

// delete file os.RemoveAll
func DelAll(path string) error{
	err:=os.RemoveAll(path)
	return err
}
