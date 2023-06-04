
layui.define(function(exports){

    layui.use(['table', 'form', 'admin', 'util'], function(){
        var table = layui.table;
        var form = layui.form;
        var admin = layui.admin;
        var util = layui.util;
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

        // 渲染 form 表单，不然搜索处的 select 不显示
        form.render();

        // 监听搜索
        form.on('submit(search-form)', function(data){
            var formData = data.field;

            //执行重载
            tableIns.reload({
                method: "post",
                url: "/api/article/search",
                page: {
                    curr: 1 //重新从第 1 页开始
                },
                where: {
                    keyword: formData.keyword
                }
            });
        });

        //生成表格
        var tableIns = table.render({
            elem: '#article-table',
            id: 'article-table-id',
            // height: 500,
            // width: 1000,
            url: '/api/article/list',
            method: 'post',
            dataType: "json",
            headers: {Token: layui.data('layuiAdmin').Token},
            contentType: 'application/json',
            page: true, //分页
            limit: 10,
            toolbar: '#article-table-toolbar', //头部盒子
            cols: [[
                {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                {field: 'name', title: 'Name', width:'23%', sort: true, fixed: 'left'},
                {field: 'link', title: 'Link'},
                {field: 'regex', title: 'Regex', hide: true},
                {field: 'rss', title: 'Rss', minWidth:520, hide: true},
                {field: 'operate', title: 'Operate', width:166, fixed: 'right', align:'center', toolbar: '#article-table-bar'}
            ]],
            done : function () {
                $('.layui-table').css("width","100%");
            }
        });

        //头工具栏事件监听
        table.on('toolbar(article-table)', function(obj){
            switch(obj.event){
                case 'refresh':
                    tableIns.reload();
                    // tableIns.reload({
                    //     page: {
                    //         curr: 1
                    //     }
                    // });
                    break;
                case 'export-all':
                    var all_count = tableIns.config.page.count;
                    if (all_count>0) {
                        var data = {
                            page: 1,
                            limit: all_count
                        };
    
                        admin.req({
                            url: '/api/article/list',
                            type: 'post',
                            dataType: "json", //期望后端返回json
                            contentType: "application/json", //发送的数据的类型
                            data: JSON.stringify(data),
                            timeout: 20000
                        }).success(function (result) {
                            if (result.code == 0){
                                table.exportFile(tableIns.config.id, result.data, 'csv');
                            }
                        });
                    } else {
                        layer.msg("当前无数据可导出", {icon: 2, time: 1000});
                    }
                    break;
            }
        });

        //行工具栏事件监听
        table.on('tool(article-table)', function(obj){
            var data = obj.data;
            switch(obj.event){
                case 'detail':
                    admin.popup({
                        type: 1, // 基本层类型：0（信息框，默认）1（页面层）2（iframe层，也就是解析content）3（加载层）4（tips层）
                        area : ["590px", '350px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        title: '查看 '+data.name,
                        content: '<div class="article-detail"></div>',
                        success: function(layero, index){
                            table.render({
                                elem: layero.find('.article-detail'),
                                width: 550,
                                data: [
                                    {x: "ID", y: data.id},
                                    {x: "Title", y: data.name},
                                    {x: "Link", y: data.link},
                                    {x: "Regex", y: data.regex},
                                    {x: "Rss", y: data.rss},
                                    {x: "Add Time", y: data.created_at},
                                ],
                                cols: [[
                                    { field: 'x', width: "18%", align:'right'},
                                    { field: 'y'}
                                ]],
                                done: function(res, curr, count){//隐藏表头
                                    layero.find('.layui-table-header').hide();
                                }
                            });
                        }
                    });
                    break;
                case 'edit':
                    admin.popup({
                        type: 1,
                        area : ["600px", '465px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        title: '编辑 '+data.name,
                        content: $('.edit-form'),
                        scrollbar: false,
                        btn: ['保存', '取消'], // 按钮：按钮1的回调是yes（也可以是btn1），而从按钮2开始，则回调为btn2: function(){}，以此类推
                        success: function(layero, index){
                            form.val('edit-form', {
                                "id": data.id,
                                "name": data.name,
                                "link": data.link,
                                "regex": data.regex,
                                "rss": data.rss,
                                "created_at": data.created_at
                            });

                            // 解决 layui 的遮罩层使用出现遮罩层覆盖弹窗情况
                            // https://article.csdn.net/h_j_c_123/article/details/104377649
                            var mask = $(".layui-layer-shade");
                            mask.appendTo(layero.parent()); //其中：layero是弹层的DOM对象

                            layero.find('.layui-layer-content').css('padding-bottom', '0px');
                            layero.find('.layui-form-label').width(61);
                            layero.find('.layui-input').width(450);
                            layero.find('.edit-form').removeClass('layui-hide');
                        },
                        yes: function(index, layero){ //更新 table 的一行
                            var formData = form.val("edit-form");
                            formData.id = parseInt(formData.id);
                            admin.req({
                                url: '/api/article/update',
                                type: 'post',
                                dataType: "json", //期望后端返回json
                                contentType: "application/json", //发送的数据的类型
                                data: JSON.stringify(formData),
                                timeout: 20000
                            }).success(function (result) {
                                if (result.code == 0){
                                    obj.update({
                                        "id": formData.id,
                                        "uname":formData.uname,
                                        "name": formData.name,
                                        "link": formData.link,
                                        "regex": formData.regex,
                                        "site_type": formData.site_type,
                                        "rss": formData.rss
                                    });
                                    layer.msg(result.msg, {icon: 1, time: 1000});
                                } else {
                                    layer.msg(result.msg, {icon: 2, time: 1000});
                                }
                            }).always(function (result) {
                                layer.close(index);
                            });
                        }
                    });
                    break;
                case 'delete':
                    layer.confirm('确定删除 '+data.name+'?', {icon: 3, shadeClose: true}, function(index){
                        data = {
                            "id": data.id,
                        }
                        admin.req({
                            url: '/api/article/delete',
                            type: 'post',
                            dataType: "json", //期望后端返回json
                            contentType: "application/json", //发送的数据的类型
                            data: JSON.stringify(data),
                            timeout: 20000
                        }).success(function (result) {
                            if (result.code == 0){
                                tableIns.reload();
                                layer.msg(result.msg, {icon: 1, time: 1000});
                            } else {
                                layer.msg(result.msg, {icon: 2, time: 1000});
                            }
                        }).always(function () {
                            layer.close(index);
                        });
                    });
                    break;
            };
        });
    });

    exports('article/site/list', {});
});