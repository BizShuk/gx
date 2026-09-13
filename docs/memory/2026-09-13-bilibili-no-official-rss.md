# 2026-09-13 — Bilibili 沒有官方 channel RSS

## 情境

vidnote 的 YouTube discover 走平台自己的
`/feeds/videos.xml?channel_id=`，不需 API key。要加 `bilibili` adapter 前，
先確認 UP 主空間有沒有對等的官方 feed，以及 `gx` 該不該像 YouTube 那樣輸出 `rss`。

## 觀察到的現象

- 空間頁沒有 RSS autodiscovery（`<link rel="alternate" type="application/rss+xml">`）。
- YouTube 同款路徑 `/feeds/videos.xml` 不存在。
- 2016 年左右的官方分區 RSS（`/rss-1.xml`、`/rss.html`）已 404。
- 未簽名的投稿 JSON API 被風控擋下；`x/web-interface/card?mid=` 仍可讀到暱稱，但那不是 feed。
- 第三方 RSSHub `/bilibili/user/video/:uid` 能合成 XML，但官方 `rsshub.app` 會被 Cloudflare challenge；少數公共實例偶爾打得開。

## 根因

Bilibili 不為空間提供 RSS。生態裡的「B 站 RSS」是打簽名 JSON API
（WBI + 反爬）再包成 XML，不是來源站的公開頁。

## 解法

- `gx youtube get channel` 輸出官方 RSS（JSON 的 `rss` 欄、`--rss`）。
- `gx bilibili get channel` 只解析 UID 與空間網址；JSON 沒有 `rss` 欄位。
- 不把 RSSHub 實例寫進輸出——公共實例不穩定，也不符合「從來源公開頁擷取」的前提。
