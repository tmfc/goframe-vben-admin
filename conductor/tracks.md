# Project Tracks

This file tracks all major tracks for the project. Each track has its own detailed plan in its respective folder.


---

## [~] Track: Modify the frontend code, we are now focusing on naive, adding a CRUD interface for the menu
*Link: [./conductor/tracks/menu_crud_naive_20251226/](./conductor/tracks/menu_crud_naive_20251226/)*

---

## [~] Track: Frontend development for web-naive, connecting remaining pages to backend APIs.
*Link: [./conductor/tracks/web_naive_frontend_20260106/](./conductor/tracks/web_naive_frontend_20260106/)*

---

## [ ] Track: 为系统增加文件上传的功能,需求如下:1.系统提供两种文件保存方式:A本地保存,B保存到S3(或兼容服务如阿里云oss,腾讯云cos或minio);2.上传过程统一如下-先调用后台上传接口(需开发)上传文件,接口返回一个路径(本地保存就返回本地的保存路径,S3保存就返回在bucket中的key),在后续业务中使用返回的路径作为需要的文件参数;规划一下后台的文件上传接口和前端的文件上传通用控件的开发
*Link: [./conductor/tracks/file_upload_20260109/](./conductor/tracks/file_upload_20260109/)*