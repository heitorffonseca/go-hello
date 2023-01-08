package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const amountOfMonitoring = 5
const delay = 5 * time.Second

func main() {
	welcome()

	for {
		app()
	}
}

func app() {
	showMenu()

	command := getCommand()

	switch command {
	case 1:
		startMonitoring()
	case 2:
		showLogs()
	case 0:
		fmt.Println("Saindo do programa")
		os.Exit(0)
	default:
		fmt.Println("Não conheço este comando!")
		os.Exit(-1)
	}
}

func welcome() {
	name := "Heitor"
	version := 1.1

	fmt.Println("Olá, sr.", name)
	fmt.Println("Este programa está na versão", version)
}

func showMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair")
}

func getCommand() int {
	var command int
	_, _ = fmt.Scan(&command)

	return command
}

func startMonitoring() {
	fmt.Println("Monitorando...")

	urls, err := readFile()

	if err != nil {
		os.Exit(-1)
	}

	for i := 0; i < amountOfMonitoring; i++ {
		for _, url := range urls {
			testUrl(url)
		}

		time.Sleep(delay)
		fmt.Println("")
	}
}

func testUrl(url string) {
	resp, _ := http.Get(url)

	if resp.StatusCode == 200 {
		fmt.Println("Site:", url, "foi carregado com sucesso!")
		registerLog(url, true)
	} else {
		fmt.Println("Site:", url, "está com problemas. Status code:", resp.StatusCode)
		registerLog(url, false)
	}
}

func readFile() ([]string, error) {
	file, err := os.Open("./sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
		return nil, err
	}

	var urls []string
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		urls = append(urls, strings.TrimSpace(line))

		if err == io.EOF {
			break
		}
	}

	_ = file.Close()

	return urls, err
}

func registerLog(url string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	_, _ = file.WriteString("Site: " + url + "\n")
	_, _ = file.WriteString("Status: " + strconv.FormatBool(status) + "\n")
	_, _ = file.WriteString("Data de teste: " + time.Now().Format("02/01/2006 15:04:05") + "\n")
	_, _ = file.WriteString("====================\n")

	_ = file.Close()
}

func showLogs() {
	fmt.Println("Exibindo Logs...")
	fmt.Println("")

	file, err := ioutil.ReadFile("./log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(file))
}
