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

$(function () {
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
    if(value == null){
        return ""
    }
    return value
}