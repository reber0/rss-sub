
layui.define(function(exports){

    layui.use(['table', 'form', 'admin'], function(){
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
                    username: formData.username,
                    email: formData.email
                }
            });
        });
        //监听搜索的 select
        form.on('select(role-select)', function(data){
            role = data.value;

            //执行重载
            tableIns.reload({
                page: {
                    curr: 1 //重新从第 1 页开始
                }
                ,where: {
                    role: role,
                }
            });
        });

        //生成表格
        var tableIns = table.render({
            elem: '#user-table',
            id: 'user-table-id',
            // height: 500,
            // width: 1000,
            url: '/api/user/list',
            method: 'post',
            dataType: "json", //期望后端返回json
            headers: {access_token: layui.data('layuiAdmin').access_token},
            contentType: 'application/json',
            page: true, //分页
            limit: 10,
            toolbar: '#table-toolbar', //头部盒子
            cols: [[
                {checkbox: true, fixed: true},
                {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                {field: 'uname', title: 'UserName', width:120, sort: true, fixed: 'left'},
                {field: 'avatar', title: 'Avatar', width:310},
                {field: 'role', title: 'Role', width:70},
                {field: 'email', title: 'Email'},
                {field: 'add_time', title: 'Add Time'},
                {field: 'operate', title: 'Operate', width:165, fixed: 'right', align:'center', toolbar: '#table-bar'},
            ]],
            done : function () {
                $('.layui-table').css("width","100%");
            }
        });

        //头工具栏事件监听
        table.on('toolbar(user-table)', function(obj){
            switch(obj.event){
                case 'refresh':
                    tableIns.reload({
                        page: {
                            curr: 1
                        }
                    });
                    break;
                case "add":
                    admin.popup({
                        type: 1,
                        title: '添加用户',
                        area: ['350px', '360px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        scrollbar: false,
                        content: $('.user-form'),
                        // content: $(".user-form").html(),
                        btn: ['添加', '取消'],
                        success: function(layero, index){
                            // 解决 layui 的遮罩层使用出现遮罩层覆盖弹窗情况
                            // https://site.csdn.net/h_j_c_123/article/details/104377649
                            var mask = $(".layui-layer-shade");
                            mask.appendTo(layero.parent()); //其中：layero是弹层的DOM对象

                            // 隐藏 id 行和 add_time 行
                            layero.find('.user-form').children().first().addClass("layui-hide");
                            layero.find('.user-form').children().last().addClass("layui-hide");

                            // // 设置 popup 的下内边距
                            layero.find('.layui-layer-content').css('padding-bottom', '0px');

                            // // 设置 label 的左内边距
                            layero.find('.user-form .layui-form-label').css('padding-left', '0px');

                            // // 显示表单
                            layero.find('.user-form').removeClass('layui-hide');

                            // 渲染 select，不然 select 不显示
                            form.render('select');
                        },
                        yes: function(index, layero){ // 添加
                            var formData = form.val("user-form");
                            admin.req({
                                url: '/api/user/add',
                                type: 'post',
                                dataType: "json", //期望后端返回json
                                contentType: "application/json", //发送的数据的类型
                                data: JSON.stringify(formData),
                                timeout: 20000,
                                success: function (result) {
                                    if (result.code == 0){
                                        tableIns.reload();
                                    }
                                    layer.msg(result.msg, {offset: '150px', icon: 1, time: 1000});
                                }
                            });
                            layer.close(index);
                        },
                        end: function(){
                            // end - 层销毁后触发的回调
                            // 还原隐藏的 id 行和 add_time 行，不然 edit 加载 html 显示不全
                            $('.user-form').children().first().removeClass("layui-hide");
                            $('.user-form').children().last().removeClass("layui-hide");

                            // 清空 form 中的 input 和 select 的值
                            $(".user-form").find("input,select").each(function(){
                                $(this).val('');
                            });
                        }
                    });
                    break;
                case 'export-select':
                    table_id = tableIns.config.id; // table-id
                    var checkStatus = table.checkStatus(table_id);
                    var data = checkStatus.data;

                    if (data.length > 0) {

                        var export_id_list = [];
                        data.forEach(function(x, i){
                            export_id_list.push(x.id);
                        });
    
                        data = {
                            page: 1,
                            limit: export_id_list.length,
                            export_id_list: export_id_list
                        }
    
                        admin.req({
                            url: '/api/user/list',
                            type: 'post',
                            dataType: "json", //期望后端返回json
                            contentType: "application/json", //发送的数据的类型
                            data: JSON.stringify(data),
                            timeout: 20000,
                            success: function (result) {
                                if (result.code == 0){
                                    table.exportFile(tableIns.config.id, result.data, 'csv');
                                }
                            }
                        });
                    } else {
                        layer.msg('没有选中数据', {offset: '150px', icon: 2, time: 1000});
                    }
                    break;
                case 'export-all':
                    all_count = tableIns.config.page.count;
                    data = {
                        page: 1,
                        limit: all_count
                    },
                    admin.req({
                        url: '/api/user/list',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(data),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                table.exportFile(tableIns.config.id, result.data, 'csv');
                            }
                        }
                    });
                    break;
                case 'delete-select':
                    // 通过 table 的唯一 id 获取选中的复选框的内容
                    site_table_id = tableIns.config.id; // table-id
                    var checkStatus = table.checkStatus(site_table_id);
                    var data = checkStatus.data;
                    if(data.length == 0) return layer.msg('未选中行', {time: 1000});

                    var id_list = [];
                    data.forEach(function(x, i){
                        id_list.push(x.id);
                    });
                    layer.confirm('确定删除 '+id_list.length+'条数据?', {icon: 3, shadeClose: true}, function(index){
                        admin.req({
                            url: '/api/user/delete',
                            type: 'post',
                            dataType: "json", //期望后端返回json
                            contentType: "application/json", //发送的数据的类型
                            data: JSON.stringify({"id_list": id_list}),
                            timeout: 20000,
                            success: function (result) {
                                if (result.code == 0){
                                    tableIns.reload();
                                }
                                layer.msg(result.msg, {offset: '150px', icon: 1, time: 1000});
                            }
                        });
                    });
                    break;
            }
        });

        //行工具栏事件监听
        table.on('tool(user-table)', function(obj){
            var data = obj.data;
            switch(obj.event){
                case 'detail':
                    admin.popup({
                        type: 1, // 基本层类型：0（信息框，默认）1（页面层）2（iframe层，也就是解析content）3（加载层）4（tips层）
                        title: '数据查看',
                        area : ["360px", '310px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        content: '<div class="detail"></div>',
                        success: function(layero, index){
                            table.render({
                                elem: layero.find('.detail'),
                                width: 320,
                                data: [
                                    {x: "ID", y: data.id},
                                    {x: "Name", y: data.uname},
                                    {x: "Role", y: data.role},
                                    {x: "email", y: data.email},
                                    {x: "Add Time", y: data.add_time},
                                ],
                                cols: [[
                                    { field: 'x', width: "30%", align:'right'},
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
                        area : ['350px', '465px'],
                        shadeClose: true, // 是否点击遮罩关闭：默认：false
                        title: '数据编辑',
                        content: $('.user-form'),
                        scrollbar: false,
                        btn: ['保存', '取消'], // 按钮：按钮1的回调是yes（也可以是btn1），而从按钮2开始，则回调为btn2: function(){}，以此类推
                        success: function(layero, index){
                            form.val('user-form', {
                                "id": data.id,
                                "uname": data.uname,
                                "password": "",
                                "role": data.role,
                                "email": data.email,
                                "add_time": data.add_time,
                            });

                            // 解决 layui 的遮罩层使用出现遮罩层覆盖弹窗情况
                            // https://article.csdn.net/h_j_c_123/article/details/104377649
                            var mask = $(".layui-layer-shade");
                            mask.appendTo(layero.parent()); //其中：layero是弹层的DOM对象

                            layero.find('.layui-layer-content').css('padding-bottom', '0px');
                            layero.find('.layui-form-label').width(61);
                            // layero.find('.layui-input').width(450);
                            layero.find('.user-form').removeClass('layui-hide');
                        },
                        yes: function(index, layero){ //更新 table 的一行
                            var formData = form.val("user-form");
                            formData.id = parseInt(formData.id);
                            admin.req({
                                url: '/api/user/update',
                                type: 'post',
                                dataType: "json", //期望后端返回json
                                contentType: "application/json", //发送的数据的类型
                                data: JSON.stringify(formData),
                                timeout: 20000
                            }).success(function (result) {
                                if (result.code == 0){
                                    obj.update({
                                        "id": formData.id,
                                        "uname": formData.uname,
                                        "role": formData.role,
                                        "email": formData.email,
                                        "add_time": formData.add_time,
                                    });
                                    layer.msg(result.msg, {icon: 1, time: 1000});
                                } else {
                                    layer.msg(result.msg, {icon: 2, time: 1000});
                                }
                            });
                            layer.close(index);
                        },
                        end: function(){
                            // 清空 form 中的 input 和 select 的值
                            $(".user-form").find("input,select").each(function(){
                                $(this).val('');
                            });
                        }
                    });
                    break;
                case 'delete':
                    layer.confirm('确定删除 '+data.username+'?', {icon: 3, shadeClose: true}, function(index){
                        var id_list = [];
                        id_list.push(data.id);
                        data = {
                            "id_list": id_list,
                        }
                        admin.req({
                            url: '/api/user/delete',
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
                        });
                    });
                    break;
            };
        });
    });

    exports('user_list_user', {});
});