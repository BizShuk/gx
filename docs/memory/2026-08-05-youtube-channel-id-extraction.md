# 2026-08-05 — YouTube 頻道 ID 不能只靠形狀比對

## 情境

`gx youtube get channel` 初版沿用手寫 curl 的思路，用 `UC[A-Za-z0-9_-]{22}`
掃整份頻道頁 HTML 取第一個匹配。

## 觀察到的現象

實測 `@YouTube` 時，每次執行得到的頻道 ID 都不同，且都不是正確的
`UCBR8-60-B28hp2BmDPdntcQ`：

```
UCRL4DE9v16a1QGF13XSmS9C
UCUqg6RAloz4i1YBLUTeZxcb
```

## 根因

頻道頁 HTML 裡有大量與頻道 ID 完全同形狀（`UC` 前綴 + 22 字元 base64url）的隨機
token——visitor data、播放清單與其他 session 相關識別碼。它們出現的位置比 canonical
link 更前面，所以「第一個匹配」抓到的是雜訊，而且每次請求都會變。

原本的 curl 寫法之所以看起來會動，是因為它比對的是
`https://www.youtube.com/channel/UC[^"]*` ——帶了網址上下文。

## 解法

`svc/youtube/extract.go` 改為依可靠度排序的多組上下文規則，一律要求 ID 出現在明確的
頻道語境中：canonical link → `channelId`/`externalId` 欄位 → channel 網址 →
`itemprop="identifier"`。ID 形狀本身只用於 `IsChannelID()` 的完整比對（已加上 `^$` 錨點）。

## 可推廣的教訓

從 HTML 抽識別碼時，**識別碼的形狀不是識別條件**，出現的上下文才是。
只要目標站的頁面夾帶同格式的 session token，形狀比對就會回傳看起來合理但錯誤的值，
且測試若只用精簡的假 HTML 不會發現——所以此類功能必須實際打一次真實端點驗證。
