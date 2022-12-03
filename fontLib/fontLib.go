// library to analyse fonts
// author: prr
// created: 1/12/2022
//
// author: prr, azul software
// date 1/12/2022
// copyright 2022 prr azul software
//

package fontLib


import (
	"os"
	"fmt"
	"encoding/binary"
)

type tableType uint32

type fileHeader struct {
	Filetype  uint32
	TablesNum uint16
	// 3 uint16 ignored
}

const (
	invalid =iota
	cmap	// character to glyph mapping
	glyf	// glyph data
	head	// font header
	hhea	// horizontal header
	hmtx	// horizontal metrics
	loca	// index to location
	maxp	// maximum profile
	name	// naming
	post	// postscript
	acnt
	ankr
	avar
	bdat
	bhed
	bloc
	bsln
//	cmap
	cvar
	cvt 
	EBSC
	fdsc
	feat
	fmtx
	fond
	fpgm
	fvar
	gasp
	gcid
//	glyf
	gvar
	hdmx
//	head
//	hhea
//	hmtx
	just
	kern
	kerx
	lcar
//	loca
	ltag
//	maxp
	meta
	mort
	morx
//	name
	opbd
	OS2
//	post
	prep
	prop
	sbix
	trak
	vhea
	vmtx
	xref
	Zapf
)

const (
	headerSize     = 4 + 4*2
	tableEntrySize = 4 * 4
)

type fontObj struct {
	FontName string
	FontType string
	Size float32
	ttfFil *os.File
	tab [48]string
	Buf *[]byte
}


func init_tbl() (tab [48]string){

	tab[0] = "invalid"
	tab[1] = "cmap"	// character to glyph mapping
	tab[2] = "glyf"	// glyph data
	tab[3] = "head"	// font header
	tab[4] = "hhea"	// horizontal header
	tab[5] = "hmtx"	// horizontal metrics
	tab[6] = "loca"	// index to location
	tab[7] = "maxp"	// maximum profile
	tab[8] = "name"	// naming
	tab[9] = "post"	// postscript
	tab[10] = "acnt"
	tab[11] = "ankr"
	tab[12] = "avar"
	tab[13] = "bdat"
	tab[14] = "bhed"
	tab[15] = "bloc"
	tab[16] = "bsln"
//	tab[8] = "cmap"
	tab[17] = "cvar"
	tab[18] = "cvt "
	tab[19] = "EBSC"
	tab[20] = "fdsc"
	tab[21] = "feat"
	tab[22] = "fmtx"
	tab[23] = "fond"
	tab[24] = "fpgm"
	tab[25] = "fvar"
	tab[26] = "gasp"
	tab[27] = "gcid"
//	tab[] = "glyf"
	tab[28] = "gvar"
	tab[29] = "hdmx"
//	tab[] = "head"
//	tab[] = "hhea"
//	tab[] = "hmtx"
	tab[30] = "just"
	tab[31] = "kern"
	tab[32] = "kerx"
	tab[33] = "lcar"
//	tab[] = "loca"
	tab[34] = "ltag"
//	tab[] = "maxp"
	tab[35] = "meta"
	tab[36] = "mort"
	tab[37] = "morx"
//	tab[] = "name"
	tab[38] = "opbd"
	tab[39] = "OS/2"
//	tab[] = "post"
	tab[40] = "prep"
	tab[41] = "prop"
	tab[42] = "sbix"
	tab[43] = "trak"
	tab[44] = "vhea"
	tab[45] = "vmtx"
	tab[46] = "xref"
	tab[47] = "Zapf"

	return tab
}

func getStrFromTtfIdx(tab [48]string, tblIdx int)(str string) {
// function that maps tabind to string
	if tblIdx<1 || tblIdx>47 {return tab[0]}
	return tab[tblIdx]
}

func getTtfIdxFromStr(tab [48] string, str string) (idx int) {
// function to return index from string
	idx = 0
	for i:=1; i< 48; i++ {
		if str == tab[i] {idx = i; break;}
	}
	return idx
}


func GetFont(fontNam string)(font *fontObj, err error) {

	var fontobj fontObj
    fontLib := "/usr/local/share/fonts/"
	ttfFil, err := os.Open(fontLib + fontNam)
	if err !=nil {return nil, fmt.Errorf("openFile %v", err)}
	defer ttfFil.Close()

	info, err := ttfFil.Stat()
	if err != nil {return nil, fmt.Errorf("Stat: %v",err)}
	size := info.Size()
	fmt.Println("size: ",size)
	buf := make([]byte, size)

	_, err = ttfFil.Read(buf)
	if err != nil {return nil, fmt.Errorf("read: %v", err)}


	// top 20 char for illustration
/*
	for i:=0; i< 20; i++ {
		fmt.Printf("%d: %d %q|", i, buf[i], buf[i])
	}
	fmt.Println()
*/
	fontobj.Buf = &buf

	// header
	// see table 4
	// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6.html
	//
	slice := buf[:4]
	tver := binary.BigEndian.Uint32(slice)
	fmt.Printf("scalar: %d\n", tver)

	slice = buf[4:6]
	numTbl := binary.BigEndian.Uint16(slice)
	fmt.Printf("table count: %d\n", numTbl)

	slice = buf[6:8]
	sr := binary.BigEndian.Uint16(slice)
	fmt.Printf("searchRange: %d\n", sr)

	slice = buf[8:10]
	es := binary.BigEndian.Uint16(slice)
	fmt.Printf("entrySel: %d\n", es)

	slice = buf[10:12]
	rs := binary.BigEndian.Uint16(slice)
	fmt.Printf("range shift: %d\n", rs)

	numTab := int(numTbl)

	// table directory
	// see table 5
	//
	fmt.Println("tbl  tag   checksum   offset  length")

	for i:=0; i< numTab; i++ {
		entry := buf[(12 +i*16): ((12 +16) +i*16)]
		tag := string(entry[:4])

		slice = entry[4:8]
		checksum := binary.BigEndian.Uint32(slice)

		slice = entry[8:12]
		offset := binary.BigEndian.Uint32(slice)

		slice = entry[12:16]
		length := binary.BigEndian.Uint32(slice)

		fmt.Printf("%3d: %4s %10d %7d %7d\n", i, tag, checksum, offset, length)
	}

	return &fontobj, nil
}
