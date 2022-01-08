# GoDiagramGenerator

[![Go Reference](https://pkg.go.dev/badge/github.com/keisuke-m123/godiagramgen.svg)](https://pkg.go.dev/github.com/keisuke-m123/godiagramgen)
[![codecov](https://codecov.io/gh/keisuke-m123/godiagramgen/branch/main/graph/badge.svg)](https://codecov.io/gh/keisuke-m123/godiagramgen)
[![ci workflow](https://github.com/keisuke-m123/godiagramgen/actions/workflows/ci.yml/badge.svg)](https://github.com/keisuke-m123/godiagramgen/actions/workflows/ci.yml)

work in progress

GoDiagramGeneratorはGo言語でのコード設計を支援するツールです。 Goコードから図を生成します。

## インストール

## 使用方法

```shell
# クラス図を生成するコマンド
godiagramgen class -h
# 使用例
godiagramgen class --recursive --output=./testingsupport/testingsupport-all-ignore-directories.puml --ignore=./testingsupport/subfolder,./testingsupport/subfolder2,./testingsupport/connectionlabels ./testingsupport

# パッケージ依存関係を生成するコマンド
godiagramgen package -h
# 使用例
godiagramgen package --output=./package-diagram.puml --theme=reddress-darkorange --ignore=./testingsupport .
```

## 生成される図

### クラス図

![クラス図](https://www.plantuml.com/plantuml/png/hLP1J-Cy4BttLypNxr9sQSHIfEu18Gehf4ehBLBP4qAL4u_18eaZspcmRF_xJfmeJXDl6wteXU2PUNv-yppojR5Csp9B9__P5ymGD7AkqPWvP_fLQPO_uyIyohnWccMGfCmOU9y0_PYrMiQbnNMYyetyXV1riflaB4DJi0I1fPAP3EsBuar9KpwzLHod09UNCFjs2Y34Sdbs9iG9sBS2GGdORcLkjrkukdA5zUyphCmwjxCJeA1RtNL1xxKgK5k99WZTnxfLnHj1-QeXueREbpz_byQEGbnnFZDWND49-E8QUAGq6IiUQuWdEYLgGUdlkoam0nIg8st0aSg8r2A95yilBQbsnsbtesW8C8N_emZsipewZcBy-5I3Eexrf-DbHlEgYq9Sl8XTBcIt71ChLPwo6DCUdtLMf9XXRiqg9cIMsBK_hcEs8MPB8GI_2BQn8l13K57dinE_B-DK9ZyDmW2_pPf3A_5TcJkgvBLOgTIxpASCWOZGU2X_TCztzEtoBUch6dNHFitdCCmu97K8SDYFozyV0NomvIguGm3N6wnWci1UqD1GuvLg74bXpQhwdjiu2bQczawfKUQYFEfUq0GtJE-bqlz_IoMK0AflvK7-4MKtmXz9z6VhYaCL7AjhgOg07Wzex14blAkVKLT_Ewl2MpmmDaqP2vLo9A5wYj2slfKhSoSnvPGK4izIwdcNYblPqw_T9mhXDhV19pd4QOZd-RuuctPzb5dQzFRYn-xYduAipkEdG8_r8ceh1rVDK_D9OUwBWSX3r6Dq7Ta1LR-RavwiPCgYKDhkgsN7MacuRjCs_RLpvujRbHMfJ9X5xaPQDqifE7fPWau9-seeJA3Adp34xiXJOzSRFZ97_a2a_stcxyllZ3nmJ21d6gniHkrHFGzLzoD1XKRwfvebi_fjOQARKbRMh8tedYfLk8asmIha91JXfqjGip2CUynyw6bObLo_Xvf6Fn1Bv3ovS9V2nUjErtJwpRx6d5Vl_3y0)

### パッケージ依存関係

![パッケージ依存関係](https://www.plantuml.com/plantuml/png/vPI_JiCm4CPtFuNfdaZ0mi3GbPadkDgdg-NuZ_ndeSgxWvGKSWCe2af9YoSd_dI_xxkJRfyBf59T9-xA4HtAX5edpBdHa6n8u0b5jiP7IE2awY1dUHBouq0foHngmHSL_AjvG_aaUk71OOwWK98fntfGmOtwhnUB9bUBZRj_U1mVkO22Da097A6V2BX8EStUvXOPJo_u5x_rLnI5YuSPZVflnFt_wKQqTCszODv_6hWytMmEhaheDIPCDplEQ6dZmwIWWaPvsf2bs84lrNMRCbKmbrhC75ExJo_jbDrEEQrocl-OhlgTV6uQkQZERY4-MheOVBmEaPVZUC-MiDlp0leZ4z7IRyZHqeN3im8peSE6MHbotXdRpS8wlduepbivNRIqNYG6sQFNrNNY7G00)

## 開発者向け

```shell
# 各pumlファイルを更新します
# これらの図はテストに使用されます
make install
make render
# テストを実行します
make test
# 静的解析
make check
```

クラス図生成について https://github.com/jfeliu007/goplantuml を参考にしています。
