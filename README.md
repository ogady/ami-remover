# ami-remover

## install

```
go install github.com/ogady/ami-remover@latest
```

## usage

引数に削除対象のAMI名（複数指定可能）を指定、-dateに削除の起点となる作成日付を指定して実行する

**AMIに紐づくEBS Snapshotも一緒に削除される。**

```sh
# Normal mode
# 2020/08/29以前に作成されたAMI（引数で指定した文字列を含むAMI）を削除する
ami-remover ${AMI NAME 1} ${AMI NAME 2} -date=20200829

# Dry run mode
ami-remover ${AMI NAME 1} ${AMI NAME 2} -date=20200829 -dry_run

```
