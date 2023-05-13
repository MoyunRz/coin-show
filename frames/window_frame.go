package frames

import (
	"coin-show/bian"
	"coin-show/cache"
	"coin-show/clients"
	"coin-show/config"
	"fmt"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"image/color"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shopspring/decimal"
)

var (
	lableItems    = map[string]*canvas.Text{}
	momeyItems    = map[string]*canvas.Text{}
	incomeItems   = &canvas.Text{}
	columnsMap    = map[string]*fyne.Container{}
	separatorMap  = map[string]*widget.Separator{}
	vcr           = &fyne.Container{}
	listContainer = &fyne.Container{}
	w             fyne.Window
)

func RunFrame() {
	// 创建应用程序
	a := app.New()
	// 创建一个窗口
	w = a.NewWindow("币价监控v0.1.8")
	// 默认配置传入
	a.Settings().SetTheme(theme.DefaultTheme())
	// 禁止调整大小
	// w.SetFixedSize(true)
	// 将窗口设置为半透明

	w.SetMaster()
	// 添加主题切换菜单项
	themeMenu := fyne.NewMenu("Theme",
		fyne.NewMenuItem("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
		fyne.NewMenuItem("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
	)
	// 创建两个标签
	listContainer = container.NewVBox()
	// 开始缓存
	LoadContainer()

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Menu",

			fyne.NewMenuItem("Coin", func() {
				// 创建一个输入框和一个按钮
				input := widget.NewEntry()
				input.SetPlaceHolder("input Coin name")

				// 创建一个输入框和一个按钮
				ibase := widget.NewEntry()
				ibase.SetPlaceHolder("input you buy price ")
				// 创建一个输入框和一个按钮
				inums := widget.NewEntry()
				inums.SetPlaceHolder("input you buy nums ")

				ncr := container.NewVBox()
				msg := canvas.NewText("", color.NRGBA{R: 255, A: 255})
				dialog := widget.NewModalPopUp(ncr, w.Canvas())

				ncr.Add(widget.NewLabel("Add Coin"))
				ncr.Add(input)
				ncr.Add(ibase)
				ncr.Add(inums)
				ncr.Add(container.NewGridWithColumns(
					2,
					widget.NewButton("Close", func() {
						msg.Text = ""
						msg.Refresh()
						dialog.Hide()
					}),
					widget.NewButton("Confirm", func() {
						if input.Text == "" {
							msg.Text = "input is not null"
							msg.Refresh()
							return
						}
						ib := "0.0"
						in := "0.0"
						if ibase.Text != "" {
							ib = ibase.Text
						}
						if inums.Text != "" {
							in = inums.Text
						}
						// 获取输入的文本
						key := strings.ToUpper(input.Text) + "USDT"
						bian.AddOrUpdateSymbolPrice(key, ib, in)
						AddContainer(strings.ToUpper(input.Text))
						msg.Text = ""
						msg.Refresh()
						dialog.Hide()
					}),
				))
				ncr.Add(msg)
				dialog.Resize(fyne.NewSize(220, 150))
				dialog.Show()
			}),
			fyne.NewMenuItem("VPN", func() {
				vpnConfig := widget.NewEntry()
				vpnConfig.SetPlaceHolder("input vpn 'http://127.0.0.1:7890'")
				ncr := container.NewVBox()
				msg := canvas.NewText("", color.NRGBA{R: 255, A: 255})
				dialog := widget.NewModalPopUp(ncr, w.Canvas())
				ncr.Add(widget.NewLabel("Revise VPN"))
				ncr.Add(vpnConfig)
				ncr.Add(container.NewGridWithColumns(
					2,
					widget.NewButton("Close", func() {
						vpnConfig.SetText("")
						msg.Text = ""
						msg.Refresh()
						dialog.Hide()
					}),
					widget.NewButton("Confirm", func() {
						if vpnConfig.Text == "" {
							msg.Text = "input is not null"
							msg.Refresh()
							return
						}
						config.Cfg.Vpn.Url = vpnConfig.Text
						bian.GetSymbol()
						vpnConfig.SetText("")
						msg.Text = ""
						msg.Refresh()
						dialog.Hide()
					}),
				))
				ncr.Add(msg)
				dialog.Resize(fyne.NewSize(220, 150))
				dialog.Show()
			}),
		),
		themeMenu,
	)
	incomeItems = canvas.NewText("0.0", color.NRGBA{R: 255, A: 255})
	totalText := widget.NewLabel("Total Income")
	totalText.TextStyle.Bold = true
	incomeItems.TextStyle.Bold = true
	w.SetMainMenu(mainMenu)
	// 将垂直框架放置在最大化容器中，并进行居中对齐
	w.Resize(fyne.NewSize(200, 300))
	vcr = container.NewVBox(
		listContainer,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, totalText, incomeItems),
	)
	// 将该容器设置为窗口的内容
	w.SetContent(vcr)
	// 显示窗口和应用程序
	w.ShowAndRun()
}

func AddContainer(key string) {
	// 获取输入的文本
	text := strings.ToUpper(key) + "USDT"
	value := clients.ReqChainLink(strings.ToUpper(key))
	if _, ok := columnsMap[text]; ok {
		return
	}

	text1 := widget.NewLabel(strings.ToUpper(key))
	text1.TextStyle.Bold = true
	text1.Alignment = fyne.TextAlignTrailing
	n, _ := strconv.ParseFloat(value, 64)
	text2 := canvas.NewText(fmt.Sprintf("%f", n), color.NRGBA{R: 128, G: 128, B: 128, A: 255})
	text2.TextStyle.Bold = true
	text2.Alignment = fyne.TextAlignCenter
	text1.Refresh()
	text2.Refresh()
	lableItems[text] = text2
	// 设置缓存
	momeyItems[text] = canvas.NewText("0.0", color.NRGBA{R: 128, G: 128, B: 128, A: 255})
	momeyItems[text].Alignment = fyne.TextAlignCenter
	text1.Refresh()
	img := container.NewGridWrap(fyne.NewSize(32, 32), LoadImage(key))
	header := container.NewHBox(container.NewCenter(img), text1)
	// 创建一个 Grid 控件
	gcls := container.NewGridWithColumns(
		3,
		header,
		container.NewCenter(lableItems[text]),
		container.NewGridWithColumns(1, momeyItems[text]),
	)

	//按钮
	btn := widget.NewButton("", func() {
		selectDo(text)
	})
	btn.Resize(fyne.NewSize(220, 40))
	columnsMap[text] = container.New(
		layout.NewMaxLayout(),
		btn,
		gcls,
	)
	listContainer.Add(columnsMap[text])
	listContainer.Refresh()
}

func LoadContainer() {
	for k, _ := range bian.NewPriceMap {
		AddContainer(strings.Split(k, "USDT")[0])
	}
}

func AsyncPrice() {
	for {
		bian.AsyncSymbolPrice()
		updateView()
		time.Sleep(time.Second * time.Duration(10))
	}
}

func updateView() {
	income := decimal.NewFromFloat(0.0)
	for k, v := range lableItems {
		if value, ok := bian.NewPriceMap[k]; ok {
			o, _ := strconv.ParseFloat(v.Text, 64)
			n, _ := strconv.ParseFloat(value.NewPrice, 64)
			price := fmt.Sprintf("%f", n)
			v.Text = price
			if o > n {
				v.Color = color.NRGBA{R: 255, A: 255}
			} else {
				v.Color = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
			}
			newPrice, _ := decimal.NewFromString(value.NewPrice)
			basePrice, _ := decimal.NewFromString(value.BasePrice)
			my := value.Nums.Mul(newPrice.Sub(basePrice))
			income = income.Add(value.Nums.Mul(newPrice.Sub(basePrice)))
			momeyItems[k].Text = my.StringFixed(2)
			momeyItems[k].Alignment = fyne.TextAlignCenter
			// 判断 decimal 值是否大于 0
			if !my.GreaterThan(decimal.NewFromInt(0)) {
				momeyItems[k].Color = color.NRGBA{R: 255, A: 255}
			} else {
				momeyItems[k].Color = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
			}

			incomeItems.Text = income.StringFixed(2)
			if !income.GreaterThan(decimal.NewFromInt(0)) {
				incomeItems.Color = color.NRGBA{R: 255, A: 255}
			} else {
				incomeItems.Color = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
			}

			v.Refresh()
			momeyItems[k].Refresh()
			incomeItems.Refresh()
		}
	}

}

func editPrice(key string) {
	ncr := container.NewVBox()
	dialog := widget.NewModalPopUp(ncr, w.Canvas())
	msg := canvas.NewText("", color.NRGBA{R: 255, A: 255})
	// 创建一个输入框和一个按钮
	ibase := widget.NewEntry()
	ibase.SetPlaceHolder("input you buy price ")
	// 创建一个输入框和一个按钮
	inums := widget.NewEntry()
	inums.SetPlaceHolder("input you buy nums ")

	ncr.Add(widget.NewLabel("Reset Price"))
	ncr.Add(ibase)
	ncr.Add(inums)
	ncr.Add(container.NewGridWithColumns(
		2,
		widget.NewButton("Close", func() {
			dialog.Hide()
		}),
		widget.NewButton("Confirm", func() {
			if ibase.Text == "" {
				msg.Text = "price is not null"
				msg.Refresh()
				return
			}
			if inums.Text == "" {
				msg.Text = "nums is not null"
				msg.Refresh()
				return
			}
			bian.AddOrUpdateSymbolPrice(key, ibase.Text, inums.Text)
			bian.GetSymbol()
			updateView()
			dialog.Hide()
		}),
	))
	ncr.Add(msg)
	dialog.Resize(fyne.NewSize(220, 150))
	dialog.Show()
}

func selectDo(text string) {
	ncr := container.NewVBox()
	dialog := widget.NewModalPopUp(ncr, w.Canvas())

	ncr.Add(container.NewBorder(
		nil,
		nil,
		widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			dialog.Hide()
		}), nil))

	ncr.Add(widget.NewLabel("Select Operate"))

	ncr.Add(container.NewGridWithColumns(
		2,
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			dialog.Hide()
			editPrice(text)
			editItemByName(text)

		}),
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			deleteItemByName(text)
			dialog.Hide()
		}),
	))
	dialog.Show()
}

func editItemByName(text string) {
	momeyItems[text].Text = bian.NewPriceMap[text].BasePrice
	momeyItems[text].Refresh()
}

func deleteItemByName(text string) {
	listContainer.Remove(columnsMap[text])
	delete(lableItems, text)
	delete(momeyItems, text)
	bian.DelCoin(text)
	listContainer.Refresh()
}

func LoadImage(key string) *canvas.Image {
	cache.WImg(fmt.Sprintf("http://124.71.12.16:7080/%s.png", strings.ToLower(key)), key)
	name := fmt.Sprintf("%s.png", key)
	path := filepath.Join(".coin_show", "icons", name)
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	img := canvas.NewImageFromURI(storage.NewFileURI(path))
	img.Resize(fyne.NewSize(20, 10))
	return img
}
