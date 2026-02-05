# aka-io - 3GPP MILENAGE AKA コマンドラインツール

Go で実装したコマンドラインツールで、3GPP MILENAGE Authentication and Key Agreement (AKA) 認証アルゴリズムの計算を行います。

## 概要

`aka-io` は 3GPP TS 35.205 に定義された MILENAGE アルゴリズムを実装し、AKA 認証の各種値を計算します。

## インストール

### ビルド

```bash
go build -o aka-io main.go
```

### 実行

```bash
./aka-io K OPc RAND AMF SQN
```

## 使用方法

### コマンド形式

```
aka-io [K] [OPc] [RAND] [AMF] [SQN]
```

### パラメータ

すべてのパラメータは Hex 表記で入力してください：

- **K**: 32 hex 文字（16 bytes）- 加入者キー
- **OPc**: 32 hex 文字（16 bytes）- K とOP から導出される値
- **RAND**: 32 hex 文字（16 bytes）- ランダムチャレンジ
- **AMF**: 4 hex 文字（2 bytes）- 認証管理フィールド
- **SQN**: 12 hex 文字（6 bytes）- シーケンス番号

### 出力

計算結果は以下の形式で標準出力に表示されます：

```
[INPUT]
K     : <input K value>
OPc   : <input OPc value>
RAND  : <input RAND value>
AMF   : <input AMF value>
SQN   : <input SQN value>
[OUTPUT]
MAC-A : <Network authentication code>
MAC-S : <Resynchronisation authentication code>
RES   : <Signed response>
CK    : <Confidentiality key>
IK    : <Integrity key>
AK    : <Anonymity key>
AKS   : <Resynchronisation anonymity key>
AUTN  : <Authentication token>
AUTS  : <Resynchronisation authentication token>
```

## 使用例

```bash
./aka-io 5C820C099FC04C3091EA0265F2824BE7 0CA072516CFCE2042AD2473BD1AAF7BD 693B2AF4B59A27391FB2D81110315811 8000 0000000003BD
```

出力例：

```
[INPUT]
K     : 5C820C099FC04C3091EA0265F2824BE7
OPc   : 0CA072516CFCE2042AD2473BD1AAF7BD
RAND  : 693B2AF4B59A27391FB2D81110315811
AMF   : 8000
SQN   : 0000000003BD
[OUTPUT]
MAC-A : 7A0C809AA9B3BEA2
MAC-S : A355B8EADF1DB880
RES   : 82DEF519E219188B
CK    : 2325A7094E6E0A60BF9477DF47578D4F
IK    : 56EDF92F645807FE8C6C9DD5A0948141
AK    : C5F7D32F0C4B
AKS   : 93124B54D00D
AUTN  : C5F7D32F0FF680007A0C809AA9B3BEA2
AUTS  : 93124B54D3B072FC66D7524CB028
```

## 実装詳細

### ライブラリ

このツールは `github.com/wmnsk/milenage v1.2.1` ライブラリを使用しています。

### AUTN と AUTS の構造

- **AUTN**: (SQN XOR AK) || AMF || MAC-A (16 bytes)
- **AUTS**: (SQN XOR AKS) || MAC-S (14 bytes)

3GPP TS 33.102 仕様に基づいた標準的な構成です。

## 依存関係

- Go 1.21+
- `github.com/wmnsk/milenage v1.2.1` - MILENAGE アルゴリズムの実装

## エラーハンドリング

以下の場合、エラーが表示されます：

- 引数の数が不足している
- Hex値が無効な形式である
- パラメータの長さが不正である

### SQN の入力形式

- SQN は 12 hex 文字の大文字・小文字を区別しないHex文字列として入力します
- 内部的には 6 bytes のバイナリ値に変換され、ビッグエンディアン（ネットワークバイト順）で処理されます
- 注：AUTN メッセージから SQN を逆算する場合は、`SQN = (SQN ⊕ AK) ⊕ AK` を使用してください

## ライセンス

MIT
