package main

import "fmt"

type Image struct {
	Name string
	ID   int64
}

type AnotherImage struct {
	Name *string
	Id   *int64
}

func main() {
	//arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	//fmt.Printf("%v\n", arr)
	//fmt.Printf("%+v\n", arr)
	//fmt.Printf("%#v\n", arr)

	//image0 := Image{
	//	Name: "name0",
	//	ID:   0,
	//}
	//
	//image1 := Image{
	//	Name: "name1",
	//	ID:   1,
	//}
	//
	//images := make([]Image, 0, 2)
	//images = append(images, image0)
	//images = append(images, image1)
	//
	//bytes, _ := json.Marshal(images)
	//fmt.Printf("修改之前：%v\n", string(bytes))
	//
	//anotherImages:= make([]AnotherImage, 0, 2)
	//for _, image := range images {
	//	anotherImage := AnotherImage{
	//		Name: &(image.Name),
	//		Id:   &image.ID,
	//	}
	//	anotherImages = append(anotherImages, anotherImage)
	//}
	//marshal, _ := json.Marshal(anotherImages)
	//fmt.Printf("重新赋值：%s", string(marshal))

	//testLenAndCap()


	testInterfaceIntV()
}

func testLenAndCap() {
	sl := make([]int64, 10)
	sl = append(sl, 1)
	println(len(sl), cap(sl)) // 11 20
	fmt.Println(sl[0], sl[10])// 0 1

	sl1 := make([]int64, 0, 10)
	fmt.Println(len(sl1), cap(sl1)) // 0 10
}

func testInterfaceIntV(){
	sl := make([]interface{}, 0,10)
	sl = append(sl, 1)
	fmt.Printf("值: %v", sl)
}

