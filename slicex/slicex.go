package slicex

// ForFunc 定义了一个带有int类型参数和bool、error类型返回值的函数类型。
// 该函数类型适用于在循环或迭代中作为停止条件判断的回调函数。
// 参数i通常代表循环的当前迭代次数或索引。
// 返回的bool值用于指示是否应该继续循环迭代。
// 返回的error值用于在迭代过程中传递错误信息，以便提前终止循环。
type ForFunc func(i int) (bool, error)

// For 执行指定次数的回调操作 代替 `for i := 0; i < count; i++`。
//
// 参数:
//
//	length 指定了回调函数应该执行的次数。
//	fn 是一个接受整数参数并返回错误的函数，它将在每次迭代中被调用，参数 i 是当前的迭代次数。
//
// 返回值:
//
//	error 如果执行过程中发生错误，则返回错误；否则返回 nil
//
// 示例:
//
//	 // 返回值 bool 为是否中途退出
//	 // 返回值 error 为是否发生错误
//	 hook := func(i int) (bool, error) {
//		fmt.Println("i = ", i)
//		if i == 3 {
//			return true, nil
//		}
//
//		return false, nil
//	 }
//
//	 if err := For(5, hook); err != nil{
//		fmt.Println(err)
//		return
//	 }
func For(length int, fn ForFunc) error {
	// 检查数组长度是否小于1，如果是，则返回长度错误。
	if length < 1 {
		return LengthError
	}
	// 检查给定的函数是否为nil，如果是，则返回函数为nil的错误。
	if fn == nil {
		return FnIsNil
	}
	// 遍历一个长度为length的空数组。
	// 这里使用EmptyArray(length)来创建一个长度符合要求的数组，但实际上不关注数组的内容，只关注索引。
	for i := range EmptyArray(length) {
		// 调用给定的函数，并传入当前的索引。
		// 函数的返回值是一个布尔值和一个错误。
		isBreak, err := fn(i)
		// 如果函数返回错误，则直接返回该错误。
		if err != nil {
			return err
		}
		// 如果函数返回true，表示应中断循环，则返回nil，表示循环已成功结束。
		if isBreak {
			return nil
		}
	}
	// 如果循环正常结束（没有返回错误，也没有提前中断），则返回nil，表示一切正常。
	return nil
}

// EmptyArray 创建一个长度为length的空结构体数组。
// 这个函数的目的是为了返回一个指定长度的切片，切片中的每个元素都是空的结构体。
// 使用空结构体的原因是，它不占用任何空间，同时提供了一个占位符的作用。
//
// 参数:
//
//	length - 指定切片的长度。
//
// 返回值:
//
//	返回一个长度为length的切片，该切片中的每个元素都是空的结构体。
//
// 使用场景:
//
//	这种空结构体切片可以用作占位符，或者在需要一个特定长度的切片但不需要具体值的场景下使用。
//
// 使用示例:
//
//	// 创建一个长度为5的空结构体切片
//	for i := range EmptyArray(5) {
//	 	fmt.Println(i)
//	}
func EmptyArray(length int) []struct{} {
	return make([]struct{}, length)
}
