/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-01-29 13:25:56
 * @LastEditTime : 2021-01-31 15:19:53
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
    
        var get_status = function(data) {
            if (data.status == 'unread'){
                return '<input type="checkbox" value="'+data.id+'" msg_type="'+data.msg_type+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch">';
            } else if (data.status == 'read') {
                return '<input type="checkbox" value="'+data.id+'" msg_type="'+data.msg_type+'" lay-skin="switch" lay-text="已读|未读" lay-filter="status-switch" checked>';
            }
        }
    
        //监听行的开关 switch 改变数据的状态(已读/未读)
        form.on('switch(status-switch)', function(obj){
            var msg_type = $(obj.elem).attr("msg_type");
            var thisTabs = tabs[msg_type];
            var update_id = obj.value;
            if (obj.elem.checked) {
                var status = "read";
            } else {
                var status = "unread";
            }
            admin.req({
                url: '/api/message/status/update',
                type: 'post',
                dataType: "json", //期望后端返回json
                contentType: "application/json", //发送的数据的类型
                data: JSON.stringify({"msg_type": msg_type, "update_id": update_id, "status": status}),
                timeout: 20000
            }).success(function (result) {
                if (result.code == 0){
                    if (result.unread_count==0){
                        $(".layui-badge-dot").addClass('layui-hide');
                    } else {
                        $(".layui-badge-dot").removeClass('layui-hide');
                    }
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
    
                if (_tabs.indexOf("system") != -1){
                    element.tabAdd("massage-tab", {
                        title: '系统消息<span class="layui-badge" id="system_num">0</span>',
                        content: '',
                        id: 'system'
                    });
    
                    // system 消息记录
                    var tableInsSystem = table.render({
                        elem: '#LAY-message-system',
                        url: '/api/message/system/list',
                        method: 'post',
                        dataType: "json",
                        headers: {access_token: layui.data('layuiAdmin').access_token},
                        contentType: 'application/json',
                        page: true, //分页
                        limit: 10,
                        cols: [[
                            {checkbox: true, fixed: true},
                            {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                            {field: 'msg_type', title: 'MsgType', width:115, sort: true,fixed: 'left'},
                            {field: 'action', title: 'Action', width:'24%'},
                            {field: 'data', title: 'Data'},
                            {field: 'add_time', title: 'Add Time', width:162},
                            {field: 'status', title: 'Status', width:95, fixed: 'right', templet:get_status},
                        ]],
                        done : function (res, curr, count) {
                            $("#system_num").html(res.unread_count);
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
                        url: '/api/message/user/list',
                        method: 'post',
                        dataType: "json",
                        headers: {access_token: layui.data('layuiAdmin').access_token},
                        contentType: 'application/json',
                        page: true, //分页
                        limit: 10,
                        cols: [[
                            {checkbox: true, fixed: true},
                            {field: 'id', title: 'ID', width:60, sort: true, fixed: 'left', align:'center'},
                            {field: 'username', title: 'UserName', width:115, sort: true,fixed: 'left'},
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
                    element.tabChange("massage-tab", "system");
                } else if (_tabs.length == 1){
                    $('.system-tab-item').remove();
                    element.tabChange("massage-tab", "user");
                }
            }
        });
    
        //区分各选项卡中的表格
        var tabs = {
            system: {
                text: '计划任务',
                id: 'LAY-message-system'
            },
            user: {
                text: '用户操作',
                id: 'LAY-message-user'
            }
        };
    
        //事件处理
        var events = {
            delete: function(othis, type){
                var thisTabs = tabs[type];
                var checkStatus = table.checkStatus(thisTabs.id);
                var data = checkStatus.data; //获得选中的数据
                if(data.length == 0) return layer.msg('未选中行', {time: 1000});
    
                var delete_id_list = [];
                data.forEach(function(x, i){
                    delete_id_list.push(x.id);
                });
                layer.confirm('确定删除选中的数据吗？', {icon: 3, shadeClose: true}, function(){
                    data = {"delete_id_list": delete_id_list, "msg_type": type}
                    admin.req({
                        url: '/api/message/delete',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(data),
                        timeout: 20000
                    }).success(function (result) {
                        if (result.code == 0){
                            if (result.unread_count==0){
                                $(".layui-badge-dot").addClass('layui-hide');
                            } else {
                                $(".layui-badge-dot").removeClass('layui-hide');
                            }
                            table.reload(thisTabs.id); //刷新表格
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        } else {
                            layer.msg(result.msg, {icon: 2, time: 1000});
                        }
                    });
                });
            },
            deleteAll: function(othis, type){
                var thisTabs = tabs[type];
                var data = {
                    "delete_all": "true",
                    "msg_type": type
                };
    
                layer.confirm('确定删除全部的数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/delete',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(data),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                if (result.unread_count==0){
                                    $(".layui-badge-dot").addClass('layui-hide');
                                } else {
                                    $(".layui-badge-dot").removeClass('layui-hide');
                                }
                                table.reload(thisTabs.id); //刷新表格
                            }
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        }
                    });
                });
            },
            read: function(othis, type){
                var thisTabs = tabs[type];
                var checkStatus = table.checkStatus(thisTabs.id);
                var data = checkStatus.data; //获得选中的数据
                if(data.length == 0) return layer.msg('未选中行', {time: 1000});
    
                var update_id_list = [];
                data.forEach(function(x, i){
                    update_id_list.push(x.id);
                });
                layer.confirm('确定已读选中的数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/status/update',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify({"update_id_list": update_id_list, "msg_type": type}),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                if (result.unread_count==0){
                                    $(".layui-badge-dot").addClass('layui-hide');
                                } else {
                                    $(".layui-badge-dot").removeClass('layui-hide');
                                }
                                table.reload(thisTabs.id); //刷新表格
                            }
                            layer.msg(result.msg, {icon: 1, time: 1000});
                        }
                    });
                });
            },
            readAll: function(othis, type){
                var thisTabs = tabs[type];
                var data = {
                    "update_all": "true",
                    "status": "read",
                    "msg_type": type
                };
    
                layer.confirm('确定已读数据吗？', {icon: 3, shadeClose: true}, function(){
                    admin.req({
                        url: '/api/message/status/update',
                        type: 'post',
                        dataType: "json", //期望后端返回json
                        contentType: "application/json", //发送的数据的类型
                        data: JSON.stringify(data),
                        timeout: 20000,
                        success: function (result) {
                            if (result.code == 0){
                                if (result.unread_count==0){
                                    $(".layui-badge-dot").addClass('layui-hide');
                                } else {
                                    $(".layui-badge-dot").removeClass('layui-hide');
                                }
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
            var type = othis.data('type'); // input 标签中的 data-type 的值 system
            var thisEvent = othis.data('events'); // input 标签中的 data-events 的值 read
            events[thisEvent] && events[thisEvent].call(this, othis, type);
        });
    });

    exports('message_list', {});
});