package main

var Icon = []string{
	":smile:",
	":laughing:",
	":smiley:",
	":sunglasses:",
	":heart:",
	":revolving_hearts:",
	":broken_heart:",
	":anger:",
	":star:",
	":star2:",
	":punch:",
	":raised_hand:",
	":clap:",
	":wedding:",
	":love_hotel:",
	":church:",
	":japan:",
	":sunrise:",
	":rainbow:",
	":triangular_flag_on_post:",
	":atm:",
	":ghost:",
	":gift:",
	":balloon:",
	":apple:",
	":green_apple:",
	":tangerine:",
	":lemon:",
	":cherries:",
	":grapes:",
	":watermelon:",
	":peach:",
	":tomato:",
	":tiger:",
	":pig:",
	":frog:",
	":chicken:",
	":beetle:",
	":fish:",
	":dolphin:",
}

//func Test2(t *testing.T) {
//	for i := 0; i < 100; i++ {
//		time.Sleep(10)
//		fmt.Println(Icon[random.RandInt(0, len(Icon)-1)])
//	}
//}
//
//func Test1(t *testing.T) {
//	fmt.Println("HelloWorld")
//	filePath := "E:\\my-study\\blog-doc\\docs\\java\\"
//	names, _ := fileutil.ListFileNames(filePath)
//	fmt.Println(names)
//	for i := range names {
//		fileLine, err := fileutil.ReadFileByLine(filePath + names[i])
//		if err != nil {
//			log.Println(err)
//		}
//		for k := range fileLine {
//			line := fileLine[k]
//			// # 一级用于作为评论来用不能增加图标
//			b, _ := regexp.MatchString("^##.*", line)
//			if b {
//				time.Sleep(30)
//				fileLine[k] = line + Icon[random.RandInt(0, len(Icon)-1)]
//				fmt.Println(fileLine[k])
//			}
//		}
//		var body string
//		for v := range fileLine {
//			body += fileLine[v] + "\n"
//		}
//		file, _ := os.OpenFile(filePath+names[i], os.O_RDWR, 0666)
//		_, err = file.WriteString(body)
//		if err != nil {
//			log.Println(err)
//		}
//		fmt.Println(fileLine)
//	}
//}
