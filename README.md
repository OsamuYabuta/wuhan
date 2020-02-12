# wuhan

これはgolang学習用に作ったサービスです。
現在開発中ですがローカルでバッチ処理からフロントエンドへの表示(トピックのみ）は確認してあります。
やってることは

- ”武漢"でのtweet検索（日本/韓国/中国/英語)、mysqlへの保存
- 形態素解析(https://github.com/OsamuYabuta/MorphologicalAnalyzerにPOSTリクエスト)、mongodbへの保存
- tfidf計算、mongodbへの保存
- トピック抽出、mysqlへの保存
- ピックアップユーザー抽出、mysqlへの保存

となります。

frontに表示する部分は途中までですが

- https://github.com/OsamuYabuta/wuhan-front

にあります。

*MorphologicalAnalyzerを動かすには別途辞書ファイルと品詞付けモデルファイルが必要です。
