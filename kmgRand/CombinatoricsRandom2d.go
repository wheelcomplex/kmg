package kmgRand

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgMath"
	"github.com/bronze1man/kmg/kmgSlice"
	"github.com/bronze1man/kmg/kmgSort"
)

/*
二维排列组合随机问题
A 组物品有若干个种类 ANumList 物品种类AKind 每个种类有ANumList[AKind]个 总数 TotalA=sum(ANumList)
B 组物品有若干个种类 BNumList 物品种类BKind 每个种类有BNumList[BKind]个 总数 TotalB=sum(BNumList)
有若干i(物品种类)和j(物品种类)的组合,有一些可以出现,有一些不允许出现
ValidCombine  ValidCombine[i][j]==true表示 i(物品种类)和j(物品种类)的组合是可以出现的
ValidCombine[i][j]==false表示 i(物品种类)和j(物品种类)的组合是不可以出现的
要求随机出min(TotalA,TotalB)个 (A,B)的组合,物品不能重复出现,要求出现顺序随机,组合随机

实现:
1.差额补全,给物品总数过小的物品组里面添加一个虚拟物品种类,将数量补全为 TotalA == TotalB,
	加入的这个虚拟物品可以和任意其他物品组合.   O(len(ANumList))
2.遍历A的物品种类,(排序数据准备)                            O(len(ANumList)*len(BNumList))
	找到这个A可行的B的物品种类的数量,
3.排序A的物品种类,(排序)                            O(len(ANumList)*log(len(ANumList)))
    下列限制分级排序:可行的B的物品种类的数量少,该A物品种类的物品数量多 排在前面
4.遍历A的物品种类的排序后顺序,(生成)                O(sum(BNumList)*sum(BNumlist))
	随机生成这个物品种类A的对应的物品种类B,将取走的具体的物品B标记,随后的随机中不取这个物品B
	4的详细步骤:
		1.将所有物品B生成一个表,包含每个B的种类,表示还没有被取走的物品B的列表
		2.遍历A的物品,计算出这个物品A可能出现的物品B的列表,
		3.从2的列表中随机取一个,保存到结果中
5.删除所有包含虚拟物品的组合
6.乱序结果列表      O(sum(BNumList))
*/
type CombinatoricsRandom2d struct {
	ANumList     []int                //key as AKindId物品种类Id,value as 该物品种类出现的数目
	BNumList     []int                //key as BKindId物品种类Id,value as 该物品种类出现的数目
	ValidCombine [][]bool             //key1 as AKindId,key2 as BKindId,value as 该种组合是否可以出现
	Output       []kmgMath.IntVector2 //key2==X as AKindId,key2==Y as BKindId ,value as (AKindId or BKindId)物品种类Id
}

func (c *CombinatoricsRandom2d) Random(r *KmgRand) (err error) {
	c.Output = []kmgMath.IntVector2{}
	aList := c.ANumList
	bList := c.BNumList
	validCombine := make([][]bool, len(c.ValidCombine))
	copy(validCombine, c.ValidCombine)
	sumA := 0
	sumB := 0
	var virtualType combinatoricsRandom2dVirtualType
	for _, num := range aList {
		sumA += num
	}
	for _, num := range bList {
		sumB += num
	}
	//1.差额补全
	switch {
	case sumA > sumB:
		virtualType = combinatoricsRandom2dVirtualTypeB
		diff := sumA - sumB
		sumB = sumA
		bList = append(bList, diff)
		for i := range aList {
			validCombine[i] = append(validCombine[i], true)
		}
	case sumA < sumB:
		virtualType = combinatoricsRandom2dVirtualTypeA
		diff := sumB - sumA
		sumA = sumB
		aList = append(aList, diff)
		thisRow := make([]bool, len(bList))
		for i := range thisRow {
			thisRow[i] = true
		}
		validCombine = append(validCombine, thisRow)
	case sumA == sumB:
		virtualType = combinatoricsRandom2dVirtualTypeNone
	}
	//fmt.Println("m0.5", bList)
	//fmt.Println("m1", virtualType)
	//2.A排序准备
	aValidBCombineNumList := make([]int, len(aList)) //key as AKindId,value as A这个种类的,可行的B的种类的数量
	for AKindId := range aList {
		for _, thisValid := range validCombine[AKindId] {
			if thisValid {
				aValidBCombineNumList[AKindId]++
			}
		}
	}
	//fmt.Println("m2", aValidBCombineNumList)
	//3.排序
	ASortOrderList := kmgSlice.IntRangeSlice(len(aList)) //key as order,value as AKindId
	kmgSort.IntLessCallbackSort(ASortOrderList, func(i int, j int) bool {
		i = ASortOrderList[i]
		j = ASortOrderList[j]
		//可行的B的物品种类的数量少
		switch {
		case aValidBCombineNumList[i] < aValidBCombineNumList[j]:
			return true
		case aValidBCombineNumList[i] > aValidBCombineNumList[j]:
			return false
		}
		//该A物品种类的物品数量多
		return aList[i] > aList[j]
	})
	//fmt.Println("m3", ASortOrderList)
	//4.生成
	theOutput := []kmgMath.IntVector2{}
	BRemainList := make([]int, sumB) //value as BKindId
	BRemainListIndex := 0
	for BKindId, num := range bList {
		for i := 0; i < num; i++ {
			BRemainList[BRemainListIndex] = BKindId
			BRemainListIndex++
		}
	}
	//fmt.Println("m3.5", BRemainList)
	for _, AKindId := range ASortOrderList {
		for i := 0; i < aList[AKindId]; i++ {
			BValidList := []int{}
			for _, BKindId := range BRemainList {
				if validCombine[AKindId][BKindId] {
					BValidList = append(BValidList, BKindId)
				}
			}
			//fmt.Println("m5.5", AKindId, BValidList, BRemainList)
			if len(BValidList) == 0 {
				return fmt.Errorf("[CombinatoricsRandom2d.Random]AKindId:%d len(BValidList)==0", AKindId)
			}
			choiceBKindId := r.ChoiceFromIntSlice(BValidList)
			//fmt.Println("m6",AKindId,choiceBKindId)
			theOutput = append(theOutput, kmgMath.IntVector2{
				X: AKindId,
				Y: choiceBKindId,
			})
			//fmt.Println("m5",BRemainList)
			kmgSlice.IntSliceRemove(&BRemainList, choiceBKindId)
		}
	}
	//fmt.Println("m4", theOutput)
	//5.删除所有包含虚拟物品的组合
	switch virtualType {
	case combinatoricsRandom2dVirtualTypeA:
		removeKindId := len(aList) - 1
		for _, row := range theOutput {
			if row.X != removeKindId {
				c.Output = append(c.Output, row)
			}
		}
	case combinatoricsRandom2dVirtualTypeB:
		removeKindId := len(bList) - 1
		for _, row := range theOutput {
			if row.Y != removeKindId {
				c.Output = append(c.Output, row)
			}
		}
	case combinatoricsRandom2dVirtualTypeNone:
		c.Output = theOutput
	}
	//6.乱序结果列表
	thisLen := len(c.Output)
	theOutput = make([]kmgMath.IntVector2, thisLen)
	permSlice := r.Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		theOutput[i] = c.Output[permSlice[i]]
	}
	c.Output = theOutput
	return
}

type combinatoricsRandom2dVirtualType int

const (
	combinatoricsRandom2dVirtualTypeA combinatoricsRandom2dVirtualType = iota
	combinatoricsRandom2dVirtualTypeB
	combinatoricsRandom2dVirtualTypeNone
)
