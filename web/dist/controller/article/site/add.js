/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-01-10 02:16:21
 * @LastEditTime: 2023-04-21 14:43:46
 */

layui.define(function(exports){

    layui.use(['form', 'admin', 'laytpl'], function(){
        var form = layui.form;
        var admin = layui.admin;
        var laytpl = layui.laytpl;
        var $ = layui.$;

        // table.render、table.reload 请求不走 admin.req，没法捕获 error， 在 401 时不会提示登录失效
        // 所以在 table.render 前添加代码捕获 complete 校验 401 弹出登录失效
        // complete 为完成请求后触发登录状态验证。即在 success 或 error 触发后触发
        $.ajaxSetup({
            complete: function (xhr, status) {
                if (xhr.status == 401) {
                    layer.msg('登录失效', {offset: '150px', icon: 2, time: 1000});
                    admin.exit();
                }
            }
        });

        // event 为 submit、lay-filter 为 rss-check 的按钮
        form.on('submit(rss-check)', function(data){
            var formData = data.field;
    
            $('#blog-msg-view').empty();
            var loading = layer.load(2);
            admin.req({
                url: '/api/article/check',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                data: JSON.stringify(formData),
                timeout: 60000
            }).success(function(result){
                if (result.code == 0){
                    var getTpl = document.getElementById('blog_msg_table').innerHTML;
                    var view = document.getElementById('blog-msg-view');
                    laytpl(getTpl).render(result, function(html){
                        view.innerHTML = html;
                    });
                }
            }).always(function(){
                layer.close(loading);
            });
        });

        form.on('submit(rss-submit)', function (data){
            var formData = data.field;

            admin.req({
                url: '/api/article/add',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                data: JSON.stringify(formData),
                timeout: 20000
            }).success(function (result) {
                if (result.code == 0){
                    var index = layer.alert(result.msg, {icon: 1, shadeClose: true});
                    layer.style(index, {width: '420px'});
                } else {
                    layer.msg(result.msg, {icon: 2, time: 1000});
                }
            }).always(function(result){
                $(":reset").click();
                $("#blog-msg-view").empty();
            });
            return false; //阻止表单跳转。如果需要表单跳转，去掉这段即可。
        });
    });

    exports('article/site/add', {});
});