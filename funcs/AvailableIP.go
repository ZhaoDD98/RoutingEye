// AvailableIP
package funcs

import (
	"fmt"
	"strconv"
)

func ToBinary(n int, bin *[8]int) { //把一个十进制数字转化为八位二进制形式，并存储在数组中
	var i = 7
	for n > 0 {
		bin[i] = n % 2
		n = n / 2
		i--
	}
}

func bin_dec(x int, n int) int {
	if n == 0 {
		return 1
	} else {
		return x * bin_dec(2, n-1)
	}
}

func Format(in string) []string { //按照格式处理输入，返回子网掩码，IP地址直接传引用
	var IP [4]int  //输入的IP地址的点分十进制表示
	var netnum int //输入的子网掩码
	//把IP,子网掩码字符串转化为数字存储在数组中
	var flag int
	var j int
	var divpos int     //子网掩码和IP地址分隔符的位置
	var binIP [32]int  //IP的二进制形式
	var binSnM [32]int //子网掩码的二进制形式
	var IPtype int     //网络类型（A,B,C）
	var occupy int
	var resthost int                //子网掩码没有占用的主机号数目
	var maxofhost int               //最大容纳的主机数
	var binNetAddress [32]int       //二进制网络地址
	var binBroadcastAddress [32]int //二进制广播地址
	var NetAddress [4]int           //十进制网络地址
	var BroadcastAddress [4]int     //十进制广播地址

	for i := 0; i < len(in); i++ {
		if in[i] == '/' {
			divpos = i
			break
		}
	}

	//字符串中提取子网掩码
	str := in[divpos+1 : len(in)]
	netnum, error := strconv.Atoi(str)
	if error == nil {
	} else {
		fmt.Println("转换错误1,", error)
	}
	//字符转中提取IP地址
	for i := 0; i < divpos; i++ {

		if in[i] == '.' {
			if flag == 0 {
				str := in[0:i]
				flag = i
				num, error := strconv.Atoi(str)
				if error == nil {
					//	fmt.Println("转换成功num/i/flag", num, i, flag)
					IP[j] = num
					j++
				} else {
					fmt.Println("转换错误1,", error)
				}

			} else {
				str := in[flag+1 : i]
				flag = i
				num, error := strconv.Atoi(str)
				if error == nil {
					//	fmt.Println("转换成功num/i/flag", num, i, flag)
					IP[j] = num
					j++
				} else {
					fmt.Println("转换错误2,", error)
				}
			}
		}
		if j == 3 {
			str := in[flag+1 : divpos]
			num, error := strconv.Atoi(str)
			if error == nil {
				//fmt.Println("转换成功num/i/flag", num, i, flag)
				IP[j] = num
				j++
			} else {
				fmt.Println("转换错误2,", error)
			}
		}
	}

	//IP地址转为32位二进制
	for i := 0; i < 4; i++ {
		var temp [8]int
		ToBinary(IP[i], &temp)

		for j := 0; j < 8; j++ {
			binIP[8*i+j] = temp[j]
		}
	}

	//子网掩码转为32位二进制
	for i := 0; i < netnum; i++ {
		binSnM[i] = 1
	}

	//判断网络类型
	if binIP[0] == 0 { //A类
		IPtype = 1
	} else if binIP[0] == 1 && binIP[1] == 0 { //B类
		IPtype = 2
	} else { //C类
		IPtype = 3
	}

	//获取子网掩码所占用的主机号数目
	if IPtype == 1 {
		occupy = netnum - 8

	} else if IPtype == 2 {
		occupy = netnum - 16

	} else if IPtype == 3 {
		occupy = netnum - 24
	}

	//获取剩余的主机号的数目
	if IPtype == 1 {
		resthost = 24 - occupy
	} else if IPtype == 2 {
		resthost = 16 - occupy
	} else {
		resthost = 8 - occupy
	}

	//fmt.Println("剩余的主机号的数目", resthost)

	//获取 最大能容纳的主机数
	maxofhost = 1
	i := resthost
	for i > 0 {
		maxofhost = maxofhost * 2
		i--
	}

	//fmt.Println("所能容纳的最大主机数为：", maxofhost)
	//fmt.Println("可用主机数为：", maxofhost-2)

	//获取网络地址
	for i := 0; i < 32; i++ {
		binNetAddress[i] = binIP[i]
	}
	j = resthost
	for i := 31; j > 0; i-- {
		binNetAddress[i] = 0
		j--
	}

	//fmt.Println("二进制网络地址为：", binNetAddress)

	//获取二进制广播地址
	for i := 0; i < 32; i++ {
		binBroadcastAddress[i] = binIP[i]
	}
	j = resthost
	for i := 31; j > 0; i-- {
		binBroadcastAddress[i] = 1
		j--
	}
	//fmt.Println("二进制广播地址为：", binBroadcastAddress)

	//二进制网络地址转换为十进制
	for i := 0; i < 8; i++ { //第一个IP（0-7位转换结果）
		if binNetAddress[i] == 1 {
			NetAddress[0] += bin_dec(2, 7-i)
		}
	}

	for i := 8; i < 16; i++ { //第二个IP（8-15位转换结果）
		if binNetAddress[i] == 1 {
			NetAddress[1] += bin_dec(2, 15-i)
		}
	}

	for i := 16; i < 24; i++ { //第三个IP（16-23位转换结果）
		if binNetAddress[i] == 1 {
			NetAddress[2] += bin_dec(2, 23-i)
		}
	}

	for i := 24; i < 32; i++ { //第四个IP（24-31位转换结果）
		if binNetAddress[i] == 1 {
			NetAddress[3] += bin_dec(2, 31-i)
		}
	}

	//二进制广播地址转换为十进制
	for i := 0; i < 8; i++ { //第一个IP（0-7位转换结果）
		if binBroadcastAddress[i] == 1 {
			BroadcastAddress[0] += bin_dec(2, 7-i)
		}
	}

	for i := 8; i < 16; i++ { //第二个IP（8-15位转换结果）
		if binBroadcastAddress[i] == 1 {
			BroadcastAddress[1] += bin_dec(2, 15-i)
		}
	}

	for i := 16; i < 24; i++ { //第三个IP（16-23位转换结果）
		if binBroadcastAddress[i] == 1 {
			BroadcastAddress[2] += bin_dec(2, 23-i)
		}
	}

	for i := 24; i < 32; i++ { //第四个IP（24-31位转换结果）
		if binBroadcastAddress[i] == 1 {
			BroadcastAddress[3] += bin_dec(2, 31-i)
		}
	}
	//fmt.Println("十进制广播地址为：", BroadcastAddress)
	//fmt.Println("十进制网络地址为：", NetAddress)

	//可用的网络地址
	//fmt.Printf("第一个可用为：%d.%d.%d.%d", NetAddress[0], NetAddress[1], NetAddress[2], NetAddress[3]+1)
	//fmt.Printf("最后一个可用为：%d.%d.%d.%d", BroadcastAddress[0], BroadcastAddress[1], BroadcastAddress[2], BroadcastAddress[3]-1)

	//输出可用列表：
	var total = maxofhost - 2
	j = 0
	flag = 32 - netnum

	ip := make([]string, total)

	//fmt.Println("可用列表：")
	if flag <= 8 {
		for i := 0; i < total; i++ {
			str1 := strconv.Itoa(NetAddress[0])
			str2 := strconv.Itoa(NetAddress[1])
			str3 := strconv.Itoa(NetAddress[2])
			str4 := strconv.Itoa(NetAddress[3] + i + 1)
			ip[i] = str1 + "." + str2 + "." + str3 + "." + str4
			//println(ip[i])
		}

	} else if flag > 8 && flag <= 16 {
		j = 0
		for i := 0; i < BroadcastAddress[2]+1; i++ {
			for j = 0; j < 256; j++ {

				str1 := strconv.Itoa(NetAddress[0])
				str2 := strconv.Itoa(NetAddress[1])
				str3 := strconv.Itoa(NetAddress[2] + i)
				if i == 0 && j == 0 {
					continue
				} else if i == BroadcastAddress[2] && j == 255 {
					break
				}
				str4 := strconv.Itoa(NetAddress[3] + j)
				ip[255*i+j] = str1 + "." + str2 + "." + str3 + "." + str4
				//fmt.Println(ip[255*i+j])
			}
		}
	} else if flag > 16 && flag <= 24 {
		j = 0
		for i := 0; i < BroadcastAddress[1]+1; i++ {
			for j = 0; j < 256; j++ {
				for k := 0; k < 256; k++ {
					str1 := strconv.Itoa(NetAddress[0])
					str2 := strconv.Itoa(NetAddress[1] + i)
					str3 := strconv.Itoa(NetAddress[2] + j)
					if i == 0 && j == 0 && k == 0 {
						continue
					} else if i == BroadcastAddress[1] && j == 255 && k == 255 {
						break
					}
					str4 := strconv.Itoa(NetAddress[3] + k)
					ip[65025*i+255*j+k] = str1 + "." + str2 + "." + str3 + "." + str4
					//fmt.Println(ip[65025*i+255*j+k])
				}
			}
		}
	} else {
		j = 0
		for i := 0; i < BroadcastAddress[0]+1; i++ {
			for j = 0; j < 256; j++ {
				for k := 0; k < 256; k++ {
					for l := 0; l < 256; l++ {
						str1 := strconv.Itoa(NetAddress[0])
						str2 := strconv.Itoa(NetAddress[1] + i)
						str3 := strconv.Itoa(NetAddress[2] + j)
						if i == 0 && j == 0 && k == 0 {
							continue
						} else if i == BroadcastAddress[1] && j == 255 && k == 255 {
							break
						}
						str4 := strconv.Itoa(NetAddress[3] + k)
						ip[16581375*i+65025*j+255*k+l] = str1 + "." + str2 + "." + str3 + "." + str4

					}
				}
			}
		}
	}
	return ip

}
