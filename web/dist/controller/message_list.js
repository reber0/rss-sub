/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-01-29 13:25:56
 * @LastEditTime: 2022-02-10 17:53:38
 */

layui.define(function(exports){

    layui.use(['admin', 'table', 'form', 'element'], function(){
        var $ = layui.$;
        var admin = layui.admin;
        var table = layui.table;
        var form = layui.form;
        var element = layui.element;
    
        element.render(); // render 后 tab 的关闭才显示
    
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

        // 检查 msg dot 状态
        var update_msg_dot = function(data) {
            admin.req({
                url: '/api/message/count',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                timeout: 20000
            }).success(function (result) {
                if (result.code == 0){
                    if (result.unread_count==0){
                        $(".layui-badge-dot").attr('class', 'layui-badge-dot layui-hide');
                    } else {
                        $(".layui-badge-dot").attr('class', 'layui-badge-dot');
                    }
                }
            });
        }
    
        var get_status = function(data) {
            var action = data.action;
            if (action.indexOf("/api/") != -1) {
                var msgtype = "api";
            } else {
                var msgtype = "user";
            }

            if (data.status == 'unread'){
                return '<input type="checkbox" value="'+data.id+'" msgtype="'+msgtype+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch">';
            } else if (data.status == 'read') {
                return '<input type="checkbox" value="'+data.id+'" msgtype="'+msgtype+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch" checked>';
            }
        }
    
        //监听行的开关 switch 改变数据的状态(已读/未读)
        form.on('switch(status-switch)', function(obj){
            var msgtype = $(obj.elem).attr("msgtype");
            var thisTabs = tabs[msgtype];

            var id_list = [];
            id_list.push(parseInt(obj.value));

            if (obj.elem.checked) {
                var status = "read";
            } else {
                var status = "unread";
            }

            var datas = {"id_list": id_list, "status": status}

            admin.req({
                url: '/api/message/update',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                data: JSON.stringify(datas),
                timeout: 20000
            }).success(function (result) {
                if (result.code == 0){
                    update_msg_dot();
                    table.reload(thisTabs.id);
                    layer.msg(result.msg, {icon: 1, time: 1000});
                } else {
                    layer.msg(result.msg, {icon: 2, time: 1000});
                }
            });
        });
    
        admin.req({
            url: '/api/message/tabs',
            type: 'post',
            dataType: "json", //期望后端返回json
            contentType: "application/json", //发送的数据的类型
            timeout: 20000,
            success: function (result) {
                var _tabs = result.data.tabs;
    
                if (_tabs.indexOf("api") != -1){
                    element.tabAdd("massage-tab", {
                        title: 'API 消息<span class="layui-badge" id="api_num">0</span>',
                        content: '',
                        id: 'api'
                    });
    
                    // api 消息记录
                    var tableInsapi = table.render({
                        elem: '#LAY-message-api',
                        url: '/api/message/api_list',
                        method: 'post',
                        dataType: "json",
                        headers: {access_token: layui.data('layuiAdmin').access_token},
                        contentType: 'application/json',
                        page: true, //分页
                        limit: 10,
                        cols: [[
                            {checkbox: true, fixed: true},
                            {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                            {field: 'username', title: 'UserName', width:120, sort: true},
                            {field: 'action', title: 'Action', width:'24%'},
                            {field: 'data', title: 'Data'},
                            {field: 'add_time', title: 'Add Time', width:162},
                            {field: 'status', title: 'Status', width:95, fixed: 'right', templet:get_status},
                        ]],
                        done : function (res, curr, count) {
                            $("#api_num").html(res.unread_count);
                            $('.layui-table').css("width","100%");
                        }
                    });
                }
                if (_tabs.indexOf("user") != -1){
                    // console.log($('#user_tab_item')[0].innerHTML);
                    element.tabAdd("massage-tab", {
                        title: '用户消息<span class="layui-badge" id="user_num">0</span>',
                        content: '',
                        id: 'user'
                    });
    
                    // user action 消息记录
                    var tableInsUser = table.render({
                        elem: '#LAY-message-user',
                        url: '/api/message/user_list',
                        method: 'post',
                        dataType: "json",
                        headers: {access_token: layui.data('layuiAdmin').access_token},
                        contentType: 'application/json',
                        page: true, //分页
                        limit: 10,
                        cols: [[
                            {checkbox: true, fixed: true},
                            {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                            {field: 'username', title: 'UserName', width:120, sort: true},
                            {field: 'action', title: 'Action', width:'24%'},
                            {field: 'data', title: 'Data'},
                            {field: 'add_time', title: 'Add Time', width:162},
                            {field: 'status', title: 'Status', width:95, fixed: 'right', templet:get_status},
                        ]],
                        done : function (res, curr, count) {
                            $("#user_num").html(res.unread_count);
                            $('.layui-table').css("width","100%");
                        }
                    });
                }
                if (_tabs.length == 2){
                    element.tabChange("massage-tab", "api");
                } else if (_tabs.length == 1){
                    $('.api-tab-item').remove();
                    element.tabChange("massage-tab", "user");
                }
            }
        });
    
        //区分各选项卡中的表格
        var tabs = {
            api: {
                text: '访问接口记录',
                id: 'LAY-message-api'
            },
            user: {
                text: '用户操作',
                id: 'LAY-message-user'
            }
        };
    
        //事件处理
        var events = {
            // 已读选中行
            read: function(othis, type){
                var thisTabs = tabs[type];
                var checkStatus = table.checkStatus(thisTabs.id);
                var data = checkStatus.data; //获得选中的数据
                if(data.length == 0) return layer.msg('未选中行', {time: 1000});
    
                var id_list = [];
                data.forEach(function(x, i){
                    id_list.push(x.id);
                });
                var datas = {"id_list": id_list, "status": "read"};

                layer.confirm('确定已读选中的数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/update',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(datas),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                update_msg_dot();
                                table.reload(thisTabs.id); //刷新表格
                            }
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        }
                    });
                });
            },
            // 已读所有行
            readAll: function(othis, type){
                var thisTabs = tabs[type];
                var datas = {
                    "msgtype": type,
                    "read_all": "true",
                };
    
                layer.confirm('确定已读数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/read_all',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(datas),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                update_msg_dot();
                                table.reload(thisTabs.id); //刷新表格
                            }
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        }
                    });
                });
            },
            // 删除选中行
            delete: function(othis, type){
                var thisTabs = tabs[type];
                var checkStatus = table.checkStatus(thisTabs.id);
                var data = checkStatus.data; //获得选中的数据
                if(data.length == 0) return layer.msg('未选中行', {time: 1000});
    
                var id_list = [];
                data.forEach(function(x, i){
                    id_list.push(x.id);
                });
                var datas = {"id_list": id_list};

                layer.confirm('确定删除选中的数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/delete',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(datas),
                        timeout: 20000
                    }).success(function (result) {
                        if (result.code == 0){
                            update_msg_dot();
                            table.reload(thisTabs.id); //刷新表格
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        } else {
                            layer.msg(result.msg, {icon: 2, time: 1000});
                        }
                    });
                });
            },
            // 删除所有行
            deleteAll: function(othis, type){
                var thisTabs = tabs[type];
                var datas = {
                    "msgtype": type,
                    "delete_all": "true"
                };
    
                layer.confirm('确定删除全部的数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/delete_all',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(datas),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                update_msg_dot();
                                table.reload(thisTabs.id); //刷新表格
                            }
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        }
                    });
                });
            }
        };
    
        $('.LAY-message-btns .layui-btn').on('click', function(){
            var othis = $(this);
            var type = othis.data('type'); // input 标签中的 data-type 的值 api
            var thisEvent = othis.data('events'); // input 标签中的 data-events 的值 read
            events[thisEvent] && events[thisEvent].call(this, othis, type);
        });
    });

    exports('message_list', {});
});
