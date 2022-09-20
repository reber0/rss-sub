
layui.define(function(exports){

    layui.use(['table', 'admin', 'form'], function(){
        var table = layui.table;
        var form = layui.form;
        var admin = layui.admin;
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

        var add_link = function(data) {
            var title = data.title;
            var url = data.url;
            return '<a class="layui-table-link" href="'+url+'" target="_blank">' + title + '</a>';
        }

        var get_status = function(data) {
            if (data.status == 'unread'){
                return '<input type="checkbox" value="'+data.id+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch">';
            } else if (data.status == 'read') {
                return '<input type="checkbox" value="'+data.id+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch" checked>';
            }
        }

        // 渲染 form 表单，不然搜索处的 select 不显示
        form.render();

        //监听搜索的提交
        form.on('submit(search-form)', function(data){
            var formData = data.field;

            //执行重载
            tableIns.reload({
                page: {
                    curr: 1 //重新从第 1 页开始
                }
                ,where: {
                    keyword: formData.keyword,
                    title: formData.title,
                    status: formData.status
                }
            });
        });
        //监听搜索的 select
        form.on('select(status-select)', function(data){
            var status = data.value;
            var keyword = $("input[name='keyword']").val();
            var title = $("input[name='title']").val();

            //执行重载
            tableIns.reload({
                page: {
                    curr: 1 //重新从第 1 页开始
                }
                ,where: {
                    keyword: keyword,
                    title: title,
                    status: status,
                }
            });
        });
        //监听行的开关 switch 改变数据的状态(已读/未读)
        form.on('switch(status-switch)', function(obj){
            // console.log(data.elem); //得到checkbox原始DOM对象
            // console.log(data.elem.checked); //是否被选中，true或者false
            // console.log(data.value); //复选框value值，也可以通过data.elem.value得到
            // console.log(data.othis); //得到美化后的DOM对象
            var status;
            if (obj.elem.checked) {
                status = "read";
            } else {
                status = "unread";
            }

            var id_list = [];
            id_list.push(parseInt(obj.value));

            admin.req({
                url: '/api/data/article/update',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                data: JSON.stringify({"id_list": id_list, "status": status}),
                timeout: 20000
            }).success(function (result) {
                if (result.code == 0){
                    tableIns.reload();
                    // tableIns.reload({
                    //     page: {
                    //         curr: 1 //重新从第 1 页开始
                    //     }
                    // });
                    layer.msg(result.msg, {icon: 1, time: 1000});
                } else {
                    layer.msg(result.msg, {icon: 2, time: 1000});
                }
            });
        });

        //生成表格
        var tableIns = table.render({
            elem: '#article-table',
            id: 'article-table-id',
            // height: 500,
            // width: 1000,
            url: '/api/data/article/list',
            method: 'post',
            dataType: 'json',
            headers: {Token: layui.data('layuiAdmin').Token},
            contentType: 'application/json',
            page: true, //分页
            limit: 10,
            where: {
                status: 'unread'
            },
            toolbar: '#article-table-toolbar', //头部盒子
            cols: [[
                {checkbox: true, fixed: true},
                {field: 'id', title: 'ID', width:65, sort: true, fixed: 'left', align:'center'},
                {field: 'name', title: 'Name', width:'15%', sort: true, fixed: 'left'},
                {field: 'title', title: 'Title', templet:add_link},
                {field: 'created_at', title: 'Created At', width:162},
                {field: 'status', title: 'Status', width:95, fixed: 'right', templet:get_status},
                {field: 'operate', title: 'Operate', width:115, fixed: 'right', align:'center', toolbar: '#article-table-bar'}
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
                case 'read-select':
                    // 通过 table 的唯一 id 获取选中的复选框的内容
                    site_table_id = tableIns.config.id; // site-table-id
                    var checkStatus = table.checkStatus(site_table_id);
                    var data = checkStatus.data;
                    if(data.length == 0) return layer.msg('未选中行', {time: 1000});

                    var id_list = [];
                    data.forEach(function(x, i){
                        id_list.push(parseInt(x.id));
                    });
                    layer.confirm('确定已读 '+id_list.length+' 条数据?', {icon: 3, shadeClose: true}, function(index){
                        admin.req({
                            url: '/api/data/article/update',
                            type: 'post',
                            dataType: "json", //期望后端返回json
                            contentType: "application/json", //发送的数据的类型
                            data: JSON.stringify({"id_list": id_list, "status": "read"}),
                            timeout: 20000
                        }).success(function (result) {
                            if (result.code == 0){
                                tableIns.reload();
                                layer.msg(result.msg, {icon: 1, time: 1000});
                            } else {
                                layer.msg(result.msg, {icon: 1, time: 1000});
                            }
                        }).always(function(){
                            layer.close(index);
                        });
                    });
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
                        area : ["590px", '310px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        title: '查看 '+data.name,
                        content: '<div class="blog-detail" ></div>',
                        success: function(layero, index){
                            table.render({
                                elem: layero.find('.blog-detail'),
                                width: 550,
                                data: [
                                    {x: "ID", y: data.id},
                                    {x: "Name", y: data.name},
                                    {x: "Title", y: data.title},
                                    {x: "URL", y: data.url},
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
                case 'delete':
                    layer.confirm('确定删除 '+data.title+'?', {icon: 3, shadeClose: true}, function(index){
                        data = {"id": data.id}
                        admin.req({
                            url: '/api/data/article/delete',
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
                        }).always(function (result) {
                            layer.close(index);
                        });
                    });
                    break;
            };
        });
    });

    exports('article/list_article', {});
});