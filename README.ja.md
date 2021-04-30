## イントロダクション

![Diagram](./img/totemo_wakariyasui_zu.jpg)

**English version is [here](./README.md)**

mtplvcapは、NikonのカメラのLive ViewをWebSocketでブラウザにリレーするマルチプラットフォーム対応 (Windows, Mac, Linux) のソフトです。

mtplvcapとOBSを組み合わせることで、NikonのカメラをHDMIキャプチャーボードなしでWebカメラにできます。お気に入りのカメラでGoogle Hangouts・Meet・Zoomを楽しみましょう！


## 動作を確認しているもの

### カメラ（確認済み）

 - Nikon D3200 (thanks @Ivisi !)
 - Nikon D3300 (thanks [@unasuke](https://github.com/unasuke) !)
 - Nikon D500 (thanks [@yasuoeto](https://github.com/yasuoeto) !)
 - Nikon D5000 (thanks [@rch850](https://github.com/rch850) !)
 - Nikon D5100 (thanks [@shigureanko](https://twitter.com/shigureanko) !)
 - Nikon D5200 (thanks [@ThatSameer](https://twitter.com/ThatSameer) !)
 - Nikon D5300
 - Nikon D5500 (thanks [@nasustim](https://github.com/nasustim) !)
 - Nikon D600 (thanks [@ohtayo](https://github.com/ohtayo) !)
 - Nikon D610 (thanks @hazlitt !)
 - Nikon D750
 - Nikon D7000 (thanks @takashi0314 !)
 - Nikon D7100 (thanks @TheMidlander !)
 - Nikon D7200 (thanks [@br_spike_love](https://twitter.com/br_spike_love) !)
 - Nikon D800E (thanks [@Higomon](https://github.com/Higomon) !)
 - Nikon D90 (thanks [@sachaos](https://github.com/sachaos) !)
 - Nikon Df
 - Nikon Z6 (thanks @ShadowXii !)
 - Nikon Z7 (thanks @zacheadams !)
 - フィードバック歓迎！お手持ちのカメラの動作可否を是非Issueにあげてください。


### カメラ（サポート外）

 - [Nikon D3000](https://github.com/puhitaku/mtplvcap/issues/18)
 - [Nikon 1 J5](https://github.com/puhitaku/mtplvcap/issues/4)


### 確認済みのOSとその他のソフト

|OS|OBS|仮想カメラ|ブラウザ|
|:-|:-|:-|:-|
|Windows 10 version 1909, Build 18363.900, amd64|25.0.8|OBS Virtualcam 2.0.5|Microsoft Edge 84.0.522.44 (Chromium Edge)|
|macOS 10.15.5 Catalina, amd64|25.0.8|obs-mac-virtualcam 3ca8f62 v1.2.0|Google Chrome 84.0.4147.89, Microsoft Edge 84.0.522.44|
|macOS 11.2.3 Big Sur, arm64|26.1.2||Microsoft Edge 88.0.705.68|
|Debian 10 Buster, amd64|25.0.7-442-ge3942061|obs-v4l2sink 1ec3c8a|Mozilla Firefox 68.10.0esr|
|Debian sid (2021-03-27), amd64|26.1.2-290-ga52012e8c|||

※ ブラウザはビデオ通話の動作テストを行ったもので、載っているブラウザでないと動作しないわけではありません


## インストール方法

注釈: このセクションにあるShellのコードスニペットは全部コピペできます。


### Windows

**重要！ Windowsでは、カメラを接続して自動でインストールされるドライバをWinUSBの汎用ドライバで置き換える必要があります。
置き換えると、手動でドライバを指定し直さない限りカメラはMTPデバイスとして認識されなくなります（カメラ自体に何かを書き込んだりはしません）
元に戻すのは至って簡単です; [こちら](#必要な場合ドライバをOS標準に戻す)のガイドを参照してください。**


#### 1. 一眼レフのドライバを置き換える

1. カメラをPCに接続します
1. Zadigを[ここ](https://zadig.akeo.ie/)からダウンロードして開きます
1. `List All Devices` にチェックを入れます

    <img alt="List All Devices にチェックを入れる" src="./img/zadig_1.png" width="400px">

1. 画面上のリストにカメラ名があるのを確認し、選択します

    <img alt="リストからカメラを選ぶ" src="./img/zadig_2.png" width="400px">

    （このスクショは置き換え後に撮ったものなので、ボタンの表記が "Reinstall Driver" となっていますが気にしないでください）

1. 中央下の選択ボックスで `WinUSB (vX.X.X.X)` を選択します
    - 他の選択肢は動作しないので注意してください

1. `Replace Driver` ボタンを押して完了するまで待ちます
    - デバイスマネージャーを起動して確認してみてください

    <img alt="インストール後のデバイスマネージャー" src="./img/devmgmt.png" width="400px">


#### 2a. ビルド済みの実行ファイルを使う

1. Windows用の実行ファイル (mtplvcap_windows_amd64.zip) を[ここ](https://github.com/puhitaku/mtplvcap/releases)からダウンロードします
1. ZIPをダブルクリックして展開します
1. `mtplvcap.exe` をダブルクリックして実行します
    - カメラのシャッターが開くのを確認してください


#### 2b. MSYS2で手でビルドする

1. MSYS2を[ここ](https://www.msys2.org/)からダウンロードしてインストールします
1. スタートメニューから "MSYS2 MSYS" を実行します
1. 依存するパッケージをインストールします

    - パッケージグループ `mingw-w64-x86_64-toolchain` とその他のパッケージを一度にインストールしようとするとエラーになる報告が上がっていて、これらは別々にインストールしたほうが良いかもしれません

	```sh
    pacman -Sy
    pacman -S mingw-w64-x86_64-toolchain
    pacman -S \
        mingw-w64-x86_64-libusb \
        mingw-w64-x86_64-go \
        mingw-w64-x86_64-pkg-config \
        git
    ```

1. PATHを追加します

    ```sh
    echo 'PATH=$PATH:/mingw64/bin:/mingw64/lib/go/bin' >> ~/.bashrc
    source ~/.bashrc
    ```

1. このリポジトリをクローンします

    ```sh
    git clone https://github.com/puhitaku/mtplvcap.git
    ```

1. `cd`してビルドして起動します

    ```sh
    cd mtplvcap
    CGO_CFLAGS='-Wno-deprecated-declarations' go build .
    ./mtplvcap.exe -debug server
    ```

    - カメラのシャッターが開くのを確認してください
    - `GOROOT=/mingw64/lib/go go build .` のようにGOROOTの指定が要るかもしれません

1. ゴール！
    - ビルドされたバイナリは移動したり再配布したりすることが可能です
    - `C:\msys64\mingw64\bin\libusb-1.0.dll` から `libusb-1.0.dll` を `mtplvcap.exe` と同じディレクトリにコピーしてくればどこでも起動できます

#### （必要な場合）ドライバをOS標準に戻す

ドライバをWinUSBに置き換えた場合、mtplvcapで使用できる代わりに、
一般的なMTPで通信するアプリ（例: 写真を取り込むアプリなど）が使用できなくなります。
以下の手順で元に戻すことができます。

1. スタートボタンを右クリックし、「デバイス マネージャー」をクリックします

    <img alt="デバイス マネージャーを起動する" src="./img/rollback_ja_1.png" width="200px">

1. "ユニバーサル シリアル バス デバイス" 以下にあるカメラを右クリックし、「ドライバーの更新」をクリックします

    <img alt="カメラを選択しドライバーの更新をクリック" src="./img/rollback_ja_2.png" width="400px">

1. 「コンピューターを参照してドライバーソフトウェアを検索」をクリックします

    <img alt="コンピューターを参照してドライバー ソフトウェアを検索をクリック" src="./img/rollback_ja_3.png" width="400px">

1. 「コンピューター上の利用可能なドライバーの一覧から選択します」をクリックします

    <img alt="コンピューター上の利用可能なドライバーの一覧から選択しますをクリック" src="./img/rollback_ja_4.png" width="400px">

1. 一覧から「MTP USB デバイス」を選択し「次へ」をクリックします

    <img alt="MTP USB デバイスを選択し次へをクリック" src="./img/rollback_ja_5.png" width="400px">

1. 「ドライバーが正常に更新されました」と表示されたら完了です

    <img alt="完了" src="./img/rollback_ja_6.png" width="400px">


### macOS

#### 1. 依存パッケージをインストールする
1. [Homebrew](https://brew.sh/)をインストールします

1. libusbをインストールします

    ```sh
    brew install libusb
    ```


#### 2a. ビルド済みの実行ファイルを使う

1. macOS用の実行ファイル (mtplvcap_macos_amd64.zip) を[ここ](https://github.com/puhitaku/mtplvcap/releases)からダウンロードします
1. ZIPを展開して起動します

    ```sh
    unzip mtplvcap_macos_amd64.zip
    ./macos/mtplvcap
    ```

    - 初回はセキュリティのアラートが出て実行できないはずです

1. 検証されていないバイナリの実行を許可する

    - 「キャンセル」でダイアログを消します (ゴミ箱には入れないでください！)

    <img alt="警告ダイアログ" src="./img/macos_warning_jp.png" width="400px">

    - 「システム環境設定」 -> 「セキュリティとプライバシー」 と開き、「このまま許可」 をクリック

    <img alt="システム環境設定" src="./img/macos_warning2_jp.png" width="400px">

    - `mtplvcap` ももう一度起動し「開く」をクリック

    <img alt="警告ダイアログ2" src="./img/macos_warning3_jp.png" width="400px">

    - カメラのシャッターが開くのを確認してください


#### 2b. 手でビルドする

1. 依存パッケージをインストールします

    ```sh
    brew install golang git
    ```

1. XCode Command Line Toolsをインストールします

    ```sh
    xcode-select --install
    ```

1. このリポジトリをクローンします

    ```sh
    git clone https://github.com/puhitaku/mtplvcap.git
    ```

1. `cd`してビルドして起動します

    ```sh
    cd mtplvcap
    CGO_CFLAGS='-Wno-deprecated-declarations' go build .
    ./mtplvcap -debug server
    ```

    - カメラのシャッターが開くのを確認してください

1. ゴール！


### Linux

ビルド済みバイナリはありますが、Linuxはディストリによって環境が大きく異なるため、手でビルドすることをおすすめします。


#### 1. 依存パッケージをインストールする

1. libusbをインストールします

    ```sh
    # For Debian/Ubuntu:
    sudo apt install libusb-1.0.0
    ```

    ```sh
    # For Arch:
    pacman --sync libusb
    ```


#### 2a. ビルド済みの実行ファイルを使う

1. Linux用の実行ファイル (mtplvcap_linux_amd64.zip) を[ここ](https://github.com/puhitaku/mtplvcap/releases)からダウンロードします

1. ZIPを展開して起動します

    ```sh
    unzip mtplvcap_linux_amd64.zip
    ./linux/mtplvcap
    ```

    - カメラのシャッターが開くのを確認してください


#### 2b. 手でビルドする

1. 依存パッケージをインストールします

    ```sh
    # For Debian/Ubuntu:
    sudo apt install golang-go libusb-1.0.0-dev
    ```

    ```sh
    # For Arch:
    pacman --sync go libusb
    ```

1. `cd`してビルドして起動します

    ```sh
    cd mtplvcap
    CGO_CFLAGS='-Wno-deprecated-declarations' go build .
    ./mtplvcap -debug server
    ```

    - カメラのシャッターが開くのを確認してください


### 使い方

```sh
$ ./mtplvcap -help
Usage of ./mtplvcap:
  -backend-go
        force gousb as libusb wrapper (not recommended)
  -debug string
        comma-separated list of debugging options: usb, data, mtp, server
  -host string
        hostname: default = localhost, specify 0.0.0.0 for public access (default "localhost")
  -port int
        port: default = 42839 (default 42839)
  -product-id string
        PID of the camera to search (in hex), default=0x0 (all) (default "0x0")
  -server-only
        serve frontend without opening a DSLR (for devevelopment)
  -vendor-id string
        VID of the camera to search (in hex), default=0x0 (all) (default "0x0")
```


#### 撮られている映像を見る

 - `http://localhost:42839/view` を開くとキャプチャされたフレームが見えます


#### ブラウザでカメラを制御する

 - `http://localhost:42839` を開くとカメラを制御するコントローラーが使えます
 - "Auto Focus" セクションは一定間隔もしくは手動でAFを動作させられます
 - "Rate Limit" セクションはフレームレートの上限を設定でき、CPU消費量の削減に使えます
 - "Information" セクションはキャプチャされているフレームの大きさ、FPS、プレビューが見えます


#### Zoom, Google Meet, Google Hangoutsなどとつなぐ

1. mtplvcapをインストールし、動作することを確認します

1. OBS (Open Broadcaster Software) を[ここ](https://obsproject.com/)からインストールします

1. OBSバーチャルカメラをインストールします（OSにより異なります; ググってください）

1. OBSの設定を開き「映像」タブを開きます

1. Live Viewの画像サイズとぴったり合うように映像サイズを設定します
    - mtplvcapを起動して `localhost:42839` を開くと「Information」セクションにサイズが載っています

    <img alt="コントローラー" src="./img/obs_1.png" width="400px">
    <img alt="解像度設定" src="./img/obs_2.png" width="400px">

1. 「ブラウザ」ソースを追加します

    <img alt="ブラウザソースを追加" src="./img/obs_3.png" width="400px">

1. 「幅」と「高さ」の値を「映像」タブの設定と同じ値に設定します

1. URLを `http://localhost:42839/view` に設定します

    <img alt="URLを設定" src="./img/obs_4.png" width="400px">

1. バーチャルカメラをONにしてチャットアプリを設定します

1. キター！！！

    <img alt="Hi!" src="./img/obs_5.png" width="400px">
    <img alt="Zoom!" src="./img/obs_6.png" width="400px">


### 諸注意

 - このソフトはアルファ版です


### 既知の問題

 - 勝手にLVが止まる
    - この自動オフは「パワーオフ時間」の設定で延長できますが、ものによって最大30分だったり、無制限だったりします
        - D5300の場合: "カスタムメニュー" -> "c AEロック・タイマー" -> "c2 パワーオフ時間" -> "カスタマイズ" -> "ライブビュー表示" -> "30分"
    - 現状では、LVが止まっていたら自動で開始するというワークアラウンドで対処しています
 - D5000: コントローラー画面で画像サイズが異常
    - 正確には640x426
 - Windowsのみ: MinTTY（MSYS2付属のターミナルエミュレータ）でmtplvcapを動かしている時、Ctrl-Cで終了すると終了処理が行われずに突然Killされる
    - 終了処理が行われないので、次の起動時に初期化に失敗したり、ケーブルの抜き差しが必要になったりします
    - これはMinTTYおよびでは知られた動作であり、mtplvcapのバグではありません
    - winptyをpacmanでインストールして、それ経由で起動すると解決します: `pacman -Sy winpty && winpty ./mtplvcap`
    - Explorerから直接mtplvcapを起動するのは問題ありません、ただしバツボタンではなく必ずCtrl-Cで終了してください


### フィードバック

 - IssueもPRも大歓迎です。[CONTRIBUTING.md](./CONTRIBUTING.md)に従って提出してください。
 - まだごく一部の機種しか動作確認できていません。是非お手持ちのカメラが動作したかどうか教えてください。何卒！


### Special Thanks

このプログラムは[github.com/hanwen/go-mtpfs](https://github.com/hanwen/go-mtpfs)からForkして大改造したものです。
go-mtpfsの成熟したMTP実装のおかげで実装のスタートを切れました。Han-Wenさん本当にありがとうございます。

[github.com/dukus/digiCamControl](https://github.com/dukus/digiCamControl)もMTPのペイロードをパースするために大いに参考にしました。
もしこのコードがなかったらmtplvcapは実装できていなかったと思います。

### ライセンス

[LICENSEファイルはこちら](./LICENSE)
