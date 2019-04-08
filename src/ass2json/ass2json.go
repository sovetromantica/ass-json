package ass2json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ScriptInfo struct {
	Title                 string
	ScriptType            string
	WarpStyle             int
	ScaledBorderAndShadow string
	YCbCr_Matrix          string
	PlayResX              int
	PlayResY              int
}

type APG struct {
	Last_Style_Storage string
	Audio_File         string
	Video_File         string
	Video_AR_Mode      int
	Video_AR_Value     float64
	Video_Zoom_Percent float64
	Scroll_Position    int
	Active_Line        int
	Video_Position     int
}

// Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
type Events struct {
	DorC    bool
	Layer   int
	Start   string
	End     string
	Style   string
	Name    string
	MarginL int
	MarginR int
	MarginV int
	Effect  string
	Text    string
}

type Styles struct {
	Name           string
	Fontname       string
	Fontsize       int
	PrimaryColor   string
	SecondaryColor string
	OutlineColor   string
	BackColor      string
	Bold           int
	Italic         int
	Underline      int
	StrikeOut      int
	ScaleX         int
	ScaleY         int
	Spacing        int
	Angle          int
	BorderStyle    int
	Outline        int
	Shadow         int
	Alignment      int
	MarginL        int
	MarginR        int
	MarginV        int
	Encoding       int
}

type EndedASS struct {
	ScriptInfo ScriptInfo
	APG        APG
	Styles     []Styles
	Events     []Events
}

type Flags struct {
	ScInfoFlag uint32
	AePrGaFlag uint32
	StylesFlag uint32
	EventFlag  uint32
}

func Ass2json(scanner *bufio.Scanner) {
	sFlags := Flags{}
	sFlags.AePrGaFlag = 0
	sFlags.EventFlag = 0
	sFlags.ScInfoFlag = 0
	sFlags.StylesFlag = 0

	sinfo := ScriptInfo{}
	sapg := APG{}
	end := EndedASS{}
	sinfomap := map[string]string{}
	aegisubgarbage := map[string]string{}

	eventsPool := [](Events){}
	stylesPool := [](Styles){}
	//scriptInfo := [](ScriptInfo){}
	//apg := [](APG){}

	for scanner.Scan() {
		// Дергаем каждый блок по отдельности
		if strings.Contains(scanner.Text(), "[Script Info]") == true {
			sFlags.ScInfoFlag = 1
		} else if strings.Contains(scanner.Text(), "[Aegisub Project Garbage]") == true {
			sFlags.ScInfoFlag = 0
			sFlags.AePrGaFlag = 1
		} else if strings.Contains(scanner.Text(), "[V4+ Styles]") == true {
			sFlags.ScInfoFlag = 0
			sFlags.AePrGaFlag = 0
			sFlags.StylesFlag = 1
		} else if strings.Contains(scanner.Text(), "[Events]") == true {
			sFlags.ScInfoFlag = 0
			sFlags.AePrGaFlag = 0
			sFlags.StylesFlag = 0
			sFlags.EventFlag = 1
		}
		// Типо если это скриптинфо из ASS
		if sFlags.ScInfoFlag == 1 {
			s := strings.Split(scanner.Text(), ":")
			if len(s) > 1 {
				if s[0][:1] == ";" {
					continue
				}
				values := s[1]
				sinfomap[s[0]] = strings.Trim(values, " ")
			}
		}
		// Мусор аеги
		if sFlags.AePrGaFlag == 1 {
			s := strings.Split(scanner.Text(), ":")
			if len(s) > 1 {
				if s[0][:1] == ";" {
					continue
				}
				values := s[1]
				aegisubgarbage[s[0]] = strings.Trim(values, " ")
			}
		}
		// Стили
		if sFlags.StylesFlag == 1 {
			styles := ParseStyle(scanner.Text())
			if len(styles.Name) <= 0 {
				continue
			}
			//log.Println(scanner.Text())
			stylesPool = append(stylesPool, styles)
		}
		// Сами субтитры
		if sFlags.EventFlag == 1 {
			event := ParseDialogueAndComments(scanner.Text())
			if len(event.Text) == 0 {
				continue
			}
			eventsPool = append(eventsPool, event)

		}
	}
	// ------- SCRIPTINFO ------
	sinfo.Title = sinfomap["Title"]
	sinfo.ScriptType = sinfomap["ScriptType"]
	warp, _ := strconv.Atoi(sinfomap["WarpStyle"])
	sinfo.WarpStyle = warp
	sinfo.ScaledBorderAndShadow = sinfomap["ScaledBorderAndShadow"]
	sinfo.YCbCr_Matrix = sinfomap["YCbCr Matrix"]
	plx, _ := strconv.Atoi(sinfomap["PlayResX"])
	ply, _ := strconv.Atoi(sinfomap["PlayResY"])
	sinfo.PlayResX = plx
	sinfo.PlayResY = ply
	// --------------------------

	// ------- ASS GARBAGE -----
	sapg.Last_Style_Storage = aegisubgarbage["Last Style Storage"]
	sapg.Audio_File = aegisubgarbage["Audio File"]
	sapg.Video_File = aegisubgarbage["Video File"]
	varmod, _ := strconv.Atoi(aegisubgarbage["Video AR Mode"])
	varmodf, _ := strconv.ParseFloat(aegisubgarbage["Video AR Value"], 4)
	vzoom, _ := strconv.ParseFloat(aegisubgarbage["Video Zoom Percent"], 4)
	scpos, _ := strconv.Atoi(aegisubgarbage["Scroll Position"])
	aline, _ := strconv.Atoi(aegisubgarbage["Active Line"])
	vpos, _ := strconv.Atoi(aegisubgarbage["Video Position"])
	sapg.Video_AR_Mode = varmod
	sapg.Video_AR_Value = varmodf
	sapg.Video_Zoom_Percent = vzoom
	sapg.Scroll_Position = scpos
	sapg.Active_Line = aline
	sapg.Video_Position = vpos
	// --------------------------
	//log.Println(sapg)
	//apg = append(apg, sapg)
	end.APG = sapg
	//end.ScriptInfo = sinfo
	//scriptInfo = append(scriptInfo, sinfo)
	end.ScriptInfo = sinfo
	end.Events = eventsPool
	end.Styles = stylesPool
	str, _ := json.Marshal(end)
	fmt.Println(string(str))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ParseStyle(sub string) Styles {
	style := Styles{}

	if (len(sub)) > 7 {
		determinator := sub[:7]
		if strings.Contains(determinator, "Style: ") {
			trimmed := strings.Trim(sub, "Style: ")
			devided := strings.Split(trimmed, ",")
			style.Name = devided[0]
			style.Fontname = devided[1]
			fsize, _ := strconv.Atoi(devided[2])
			style.Fontsize = fsize
			style.PrimaryColor = devided[3]
			style.SecondaryColor = devided[4]
			style.OutlineColor = devided[5]
			style.BackColor = devided[6]
			bold, _ := strconv.Atoi(devided[7])
			italic, _ := strconv.Atoi(devided[8])
			underline, _ := strconv.Atoi(devided[9])
			strikeout, _ := strconv.Atoi(devided[10])
			scalex, _ := strconv.Atoi(devided[11])
			scaley, _ := strconv.Atoi(devided[12])
			spacing, _ := strconv.Atoi(devided[13])
			angle, _ := strconv.Atoi(devided[14])
			borderstyle, _ := strconv.Atoi(devided[15])
			outline, _ := strconv.Atoi(devided[16])
			shadow, _ := strconv.Atoi(devided[17])
			aligment, _ := strconv.Atoi(devided[18])
			marginl, _ := strconv.Atoi(devided[19])
			marginr, _ := strconv.Atoi(devided[20])
			marginv, _ := strconv.Atoi(devided[21])
			encoding, _ := strconv.Atoi(devided[22])
			style.Bold = bold
			style.Italic = italic
			style.Underline = underline
			style.StrikeOut = strikeout
			style.ScaleX = scalex
			style.ScaleY = scaley
			style.Spacing = spacing
			style.Angle = angle
			style.BorderStyle = borderstyle
			style.Outline = outline
			style.Shadow = shadow
			style.Alignment = aligment
			style.MarginL = marginl
			style.MarginR = marginr
			style.MarginV = marginv
			style.Encoding = encoding
			return style
		}
	}
	return style
}

func ParseDialogueAndComments(sub string) Events {
	event := Events{}

	// Да да мы все знаем даблокод, не ну а хуле
	if len(sub) > 9 {
		determinator := sub[:9]
		if strings.Contains(determinator, "Dialogue:") == true {
			trimmed := strings.Trim(sub, "Dialogue: ")
			devided := strings.Split(trimmed, ",")

			mul := strings.Join(devided[9:], ",")
			// Inputting Dialogue event
			// Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
			event.DorC = true
			layer, _ := strconv.Atoi(devided[0])
			marginl, _ := strconv.Atoi(devided[5])
			marginr, _ := strconv.Atoi(devided[6])
			marginv, _ := strconv.Atoi(devided[7])
			event.Layer = layer
			event.Start = devided[1]
			event.End = devided[2]
			event.Style = devided[3]
			event.Name = devided[4]
			event.MarginL = marginl
			event.MarginR = marginr
			event.MarginV = marginv
			event.Effect = devided[8]
			event.Text = mul
		} else if strings.Contains(determinator, "Comment:") == true {
			trimmed := strings.Trim(sub, "Dialogue: ")
			devided := strings.Split(trimmed, ",")

			mul := strings.Join(devided[9:], ",")
			// Inputting Dialogue event
			// Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
			event.DorC = false
			layer, _ := strconv.Atoi(devided[0])
			marginl, _ := strconv.Atoi(devided[5])
			marginr, _ := strconv.Atoi(devided[6])
			marginv, _ := strconv.Atoi(devided[7])
			event.Layer = layer
			event.Start = devided[1]
			event.End = devided[2]
			event.Style = devided[3]
			event.Name = devided[4]
			event.MarginL = marginl
			event.MarginR = marginr
			event.MarginV = marginv
			event.Effect = devided[8]
			event.Text = mul
		}
	}

	return event
}
