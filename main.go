package main

import (
	"os"
	"fmt"
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Columns struct {
	table_name string
	table_name_alias string
	column_name string
	column_name_alias string
	column_id int
	srccols map[string]string
}


func getTableName(sql string)string{
	reg_insert:=regexp.MustCompile("insert into ([0-9a-zA-Z._]+)")
	reg_join:=regexp.MustCompile("left|right|inner ([0-9a-zA-Z._]+)")
	reg_from:=regexp.MustCompile("from ([0-9a-zA-Z._]+)")
	res:=reg_insert.ReplaceAllString(sql,"$1")
	if len(res) > 0 {
		return res
	}
	res=reg_join.ReplaceAllString(sql,"$1")
	if len(res) > 0 {
		return res
	}
	res=reg_from.ReplaceAllString(sql,"$1")
	if len(res) > 0 {
		return res
	}else {
		return ""
	}

}

func MatchPairs(str string,bdelim string,edelim string)([]int,[]int,bool){
	cnt:=0

	reg1:=regexp.MustCompile(bdelim)
	bidx:=reg1.FindAllStringIndex(str,-1)

	reg2:=regexp.MustCompile(edelim)
	eidx:=reg2.FindAllStringIndex(str,-1)

	i,j:=0,0
	lenb := len(bidx)
	lene := len(eidx)

	if lenb < 1 && lene < 1{
		return nil,nil,false
	}

	for {
		if bidx[i][0] < eidx[j][0] {
			cnt++

			if i+1 >= lenb{
				break
			}else{
				i++
			}
		}else{
			break
		}
	}

	if cnt > lene {
		return nil,nil,false
	}else{
		return bidx[0],eidx[cnt - 1],true
	}

}

func GetColumns(str string)[]string{
	cols := strings.Split(str,",")
	for i := range cols{
		if strings.Contains(cols[i],"CASE") == false && strings.Contains(cols[i],"(") == false{
			fmt.Println(cols[i])
			continue
		}
		casecnt := strings.Contains(cols[i],"CASE")
		khcnt := strings.Contains(cols[i],"(")

		if strings.Contains(cols[i],"CASE") == true{
			strings.Count()
		}

	}


	return nil
}

func main(){
	fr,err := os.Open("D:\\我接收到的文件\\testplsqlreader.txt")
	if err != nil {
		fmt.Println("open file error")
		return
	}

	defer fr.Close()
	var allsql string
	br := bufio.NewReader(fr)
	iscombg := false
	for {
		sql,err := br.ReadString('\n')
		if err == io.EOF{
			break
		}

		//delete blank

		regblank := regexp.MustCompile("^[ ]+|[ ]+$")
		sql = regblank.ReplaceAllString(sql,"")

		regcom := regexp.MustCompile("/\\*.*\\*/")
		sql = regcom.ReplaceAllString(sql,"")

		//delete single row comment
		reg1 := regexp.MustCompile("--.*\n")
		sql = reg1.ReplaceAllString(sql,"")

		//delete multiple rows comment
		if !iscombg {
			reg2 := regexp.MustCompile("/\\*.*\n")
			idx := reg2.FindStringIndex(sql)
			if idx != nil {
				iscombg = true
				sql = reg2.ReplaceAllString(sql, "")
				if len(sql) > 0 {
					sql = strings.Replace(sql,"\r\n","",1)
					allsql += " "+sql
				}
				continue
			}
		}else{
			reg3 := regexp.MustCompile("^.*\\*/")
			idx := reg3.FindStringIndex(sql)
			if idx != nil {
				sql = reg3.ReplaceAllString(sql, "")
				iscombg = false
				if len(sql) > 0 {
					sql = strings.Replace(sql,"\r\n","",1)
					allsql += " "+sql
					continue
				}
			}else {
				continue
			}
		}

		sql = strings.Replace(sql,"\r\n","",1)

		sql = strings.ToUpper(sql)
		if len(sql) > 0 {
			allsql += " "+sql
		}


	}
	regblank := regexp.MustCompile("([ ]*,[ ]*)")
	allsql = regblank.ReplaceAllString(allsql,",")
	fmt.Println(allsql)
	bpos,epos,_:= MatchPairs(allsql,"INTO","\\(")
	fmt.Println("table_name:",allsql[bpos[1]:epos[0]-1])
	//bg:=bpos[1]
	//ed:=epos[0]
	//lastpos := epos[0] - 1
	bpos,epos,_= MatchPairs(allsql,"\\(","\\)")
	//fmt.Println(allsql[bpos[1]:epos[0]-1])
	GetColumns(allsql[bpos[1]:epos[0]-1])
	fmt.Println("-------------------")
	bpos,epos,_= MatchPairs(allsql,"SELECT","FROM")
	//fmt.Println(allsql[bpos[1]:epos[0]-1])
	GetColumns(allsql[bpos[1]:epos[0]-1])


}
