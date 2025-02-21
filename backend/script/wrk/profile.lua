-- 设置请求头
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQxNDQwMTksIlVJRCI6MywiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4zOS4wIn0.hPb1Ux2NQ7SHIQzJYC5aqcq4JaVr5vnp57bc4EGscLk"
wrk.method = "GET"

-- 定义请求函数
request = function()
    -- 构造请求体
    --local body = '{"name": "test", "age": 30}'
    return wrk.format(wrk.method, "/users/profile", wrk.headers)
    --print("Request URL:", wrk.format("GET", "/users/profile?email=cc@163.com", wrk.headers))
end