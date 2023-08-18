package mq

import "fmt"

// MessageHandler 定义消息处理函数
// 按照注释的要求去定义消息处理函数：return的值是一个函数，此函数是在需要处理的 service 层内部的函数
// 例如：return func(msg string) { service.AddComment(...) }
// 可以选择一个返回值，也可以选择两个返回值，第一个返回值是一个函数，第二个返回值是一个 error
type MessageHandler func(string)

func AddComment(msg string) {
	fmt.Println("Adding comment:", msg)
}

func DeleteComment(msg string) {
	fmt.Println("Deleting comment:", msg)
}

func AddLike(msg string) {
	fmt.Println("Adding like:", msg)
}

func RemoveLike(msg string) {
	fmt.Println("Removing like:", msg)
}

func AddFollow(msg string) {
	fmt.Println("Adding follow:", msg)
}

func RemoveFollow(msg string) {
	fmt.Println("Removing follow:", msg)
}
