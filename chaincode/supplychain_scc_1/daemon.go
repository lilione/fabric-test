package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func handOffItemToNextProvider(
	idxInputProvider string,
	maskedInputProvider string,
	idxOutputProvider string,
	maskedOutputProvider string,
	idxAmt string,
	maskedAmt string,
	itemID string,
	prevSeq string,
	seq string,
	sharePrevOutputProvider string,
	sharePrevAmt string) {

	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v1/hand_off_item_to_next_provider.py", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	//cmd := exec.Command("python3", "apps/fabric/src/supplychain/v1/hand_off_item_to_next_provider.py", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
	//cmd.Dir = "/opt/gopath/src/github.com/lilione/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/log.txt")
	//file, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	fmt.Println(cmd)
	fmt.Println("starting cmd")
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}
	fmt.Println("cmd started")
}

//func main() {
//	handOffItemToNextProvider(
//		"12",
//		"31735365036769119937719688508119359025704093086069425564197049281662336681743",
//		"13",
//		"31642726331775536829915739815154062633111360601414380494429154729856685199456",
//		"14",
//		"49552563574017699443549457631995460639422346620873144610294471217117283829075",
//		"2",
//		"0",
//		"1",
//		"{50259279881865192434443071006659689597269766008958498194490654843566364568305}",
//		"{45651947582898415568088281581893079250520598662323323056668936736811008813408}")
//}