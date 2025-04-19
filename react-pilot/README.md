# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

```shell
npm create vite@4

npm i
npm i eslint vite-plugin-eslint eslint-config-react-app --save-dev

# then create file named .eslintrc.json

# added the json:
{
  "extends": "react-app"
}

# in the vite.config.js, and add eslint():
import eslint from "vite-plugin-eslint";

npm i styled-components

npm i react-router-dom@6

# 20250301
npm i react-icons

# 使用 supabase 作为后端，包含数据库 api 和图片存储等
https://supabase.com/dashboard/project/nvbtdjgdbhgsgccuwkap/editor/29291?schema=public

npm install --save @supabase/supabase-js

# request query calling
npm i @tanstack/react-query@4
npm i @tanstack/react-quary-devtools@4

npm i react-hot-toast
npm i react-hook-form@7

npm install date-fns@latest
npm cache clean --force


# 使用这个地址做邮箱确认
https://temp-mail.org/

npm i recharts@2
npm i react-error-boundary

# 构建项目
npm run build

netlify.com
https://the-wild-oasis-half-coconut.netlify.app/login
```
