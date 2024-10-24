// 仿java
String.prototype.startWith = function (str) {
    if (str == null || str == "" || this.length == 0
        || str.length > this.length)
        return false;
    if (this.substr(0, str.length) == str)
        return true;
    else
        return false;
    return true;
};

// 仿java
String.prototype.endWith = function (str) {
    if (str == null || str == "" || this.length == 0
        || str.length > this.length)
        return false;
    if (this.substring(this.length - str.length) == str)
        return true;
    else
        return false;
    return true;
};

Date.prototype.format = function (pattern = "yyyy-MM-dd hh:mm:ss") {
    const o = {
        "M+": this.getMonth() + 1, // month
        "d+": this.getDate(), // day
        "h+": this.getHours(), // hour
        "m+": this.getMinutes(), // minute
        "s+": this.getSeconds(), // second
        "q+": Math.floor((this.getMonth() + 3) / 3), // quarter
        "S": this.getMilliseconds()
        // millisecond
    };

    if (/(y+)/.test(pattern)) {
        pattern = pattern.replace(RegExp.$1, (this.getFullYear() + "")
            .substr(4 - RegExp.$1.length));
    }

    for (var k in o) {
        if (new RegExp("(" + k + ")").test(pattern)) {
            pattern = pattern.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k] :
                ("00" + o[k]).substr(("" + o[k]).length));
        }
    }
    return pattern;
}

$(function () {
    initDropdown()
    if ($(".navbar").length > 0) {
        initTopBar();
    }
});

/**
 * ajax请求失败通用处理
 * @param res
 */
const $ajaxFail = function (res) {
    const data = res.responseJSON
    if (data.code === 1) {
        window.location.href = "/admin/login"
    } else {
        alert(data.msg)
    }
}

/**
 * 退出登录
 */
function logout() {
    $.ajax({
        url: "/admin/login/logout",
        type: "POST"
    }).fail($ajaxFail).done((data) => {
        window.location.href = "/admin/login"
    })
}

/**
 * 重置账户
 */
function reinit() {
    $.ajax({
        url: "/admin/index/reinit",
        type: "POST"
    }).fail($ajaxFail).done((data) => {
        window.location.href = "/admin/login"
    })
}

function getParam(key) {

    // 获取当前页面的 URL
    const urlParams = new URLSearchParams(window.location.search);

    // 获取单个参数值
    const value = urlParams.get(key);
    if (value == null) {
        return ""
    }
    return value
}

/**
 * 初始化dropdown数据
 */
function initDropdown(){
    let query = ""
    $.each($("select[dropdown-tag]"),(_,obj)=>{
        query += $(obj).attr("dropdown-tag") + "=1&"
    })
    if(query === ""){
        return
    }
    $.ajaxByData("common/dropdown?" + query).success(data=>{
        console.log(data)
        for(const key in data){
            const $select = $(`select[dropdown-tag="${key}"]`)
            const options = data[key]
            options.forEach(item=>{
                $select.append(`<option value="${item.Value}">${item.Label}</option>`)
            })
        }
    }).post()
}