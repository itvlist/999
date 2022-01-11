/*!
 * Stui v3.0 Copyright 2016-2018 http://v.shoutu.cn
 * Email 726662013@qq.com,admin@shoutu.cn
 */
var stui = {
	'browser': {//浏览器
		url: document.URL,
		domain: document.domain,
		title: document.title,
		language: (navigator.browserLanguage || navigator.language).toLowerCase(),
		canvas: function() {
			return !!document.createElement("canvas").getContext
		}(),
		useragent: function() {
			var a = navigator.userAgent;
			return {
				mobile: !! a.match(/AppleWebKit.*Mobile.*/),
				ios: !! a.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/),
				android: -1 < a.indexOf("Android") || -1 < a.indexOf("Linux"),
				iPhone: -1 < a.indexOf("iPhone") || -1 < a.indexOf("Mac"),
				iPad: -1 < a.indexOf("iPad"),
				trident: -1 < a.indexOf("Trident"),
				presto: -1 < a.indexOf("Presto"),
				webKit: -1 < a.indexOf("AppleWebKit"),
				gecko: -1 < a.indexOf("Gecko") && -1 == a.indexOf("KHTML"),
				weixin: -1 < a.indexOf("MicroMessenger")
			}
		}()
	},
	'images': {//图片处理
		'lazyload': function() {
			$(".lazyload").lazyload({
				effect: "fadeIn",
				threshold: 200,
				failurelimit: 15,
				skip_invisible: false
			})
		},
		'qrcode': function() {
			$("img.qrcode").attr("src", "https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=" + encodeURIComponent(stui.browser.url) + "");
		}
	},
	'common': {//公共基础
		'tab': function() {
			$(".tab li").on('click',function(){
			    $(".tab li.active").removeClass('active');
			    $(this).addClass('active');
			    var index = $(this).index();
				$(".tab-content .item").eq(index).addClass('active').siblings().removeClass('active');
	  		});
		  	$(".down-tab li").on('click',function(){
			    $(".down-tab li.active").removeClass('active');
			    $(this).addClass('active');
			    var index = $(this).index();
			    $(this).parent().parent().find("h3").html($(".down-tab li.active").html());
				$(".down-content .down-item").eq(index).addClass('active').siblings().removeClass('active');
				$(".down-tab").hide();			
	  		});
		  	$(".play-tab li").on('click',function(){
			    $(".play-tab li.active").removeClass('active');
			    $(this).addClass('active');
			    var index = $(this).index();
			    $(this).parent().parent().find("h3").html($(".play-tab li.active").html())
				$(".play-content .play-item").eq(index).addClass('active').siblings().removeClass('active');
				$(".play-tab").hide();			
	  		});
		  	$(".play-switch").on('click',function(){
		  		$(".play-tab").toggle();
		  	});
		  	$(".down-switch").on('click',function(){
		  		$(".down-tab").toggle();
		  	});
		},	
		'history': function() {
			if($.cookie("recente")){
			    var json=eval("("+$.cookie("recente")+")");
			    var list="";
			    for(i=0;i<json.length;i++){
			        list = list + "<li class='top-line'><a href='"+json[i].vod_url+"' title='"+json[i].vod_name+"'><span class='pull-right text-red'>"+json[i].vod_part+"</span>"+json[i].vod_name+"</a></li>";
			    }
			    $("#stui_history").append(list);
			}
			else
	            $("#stui_history").append("<p style='padding: 80px 0; text-align: center'>您还没有看过影片哦</p>");
		   
		    $(".historyclean").on("click",function(){
		    	$.cookie("recente",null,{expires:-1,path: '/'});
		    })		    
		},
		'collapse': function() {
			$("a.detail-more").on("click",function(){
				$(this).parent().find(".detail-sketch").addClass("hide");
				$(this).parent().find(".detail-content").css("display","");
				$(this).remove();
			})
		},
		'more': function() {
			$(".menu-switch").on('click',function(){
		  		var display = PlaySide.css('display');
		  		if(display == 'block'){
		  			PlaySide.hide(); 
		  			PlayLeft.css("width","100%");
		  			$(this).find("span").html("打开菜单");
				}else{
					PlaySide.show();  
					PlayLeft.css("width","75%");
					$(this).find("span").html("关闭菜单");
				}
			})
		  	
		  	$(".open-desc").on('click',function(){
		  		$(".data-more").slideToggle("slow")
		  	})
		  	
		  	var date = new Date;
			var h = date.getHours();  //时
			var minute = date.getMinutes()  //分
			if(h<10){
				h = "0"+h;
			}
			if(minute<10){
				minute = "0"+minute;
			}
			$(".date").html('<span>'+h+":"+minute+"</span>");
		}
	}	
};
$(document).ready(function() {	
	stui.images.lazyload();
	stui.images.qrcode();
	stui.common.tab();
	stui.common.history();
	stui.common.collapse();
	stui.common.more();
        $('#alBRyGHxuDYU span').click(function() {
            $('#alBRyGHxuDYU').hide()
            sessionStorage.setItem('isShowL', 'hide')

        })
        $('#arBRyGHxuDYU span').click(function() {
            $('#arBRyGHxuDYU').hide()
            sessionStorage.setItem('isShowR', 'hide')
        })
        $('.advertise span').click(function() {
            $('.advertise').hide()
            sessionStorage.setItem('isShowH5', 'hide')
        })
        // 右边广告
        function myFnR() {
            let ret = sessionStorage.getItem('isShowR')
            if (ret == 'hide') {
                $('#arBRyGHxuDYU').hide()
            } else {
                $('#arBRyGHxuDYU').show()
            }
        }
        myFnR()
        // 左边广告
        function myFnL() {
            let ret = sessionStorage.getItem('isShowL')
            if (ret == 'hide') {
                $('#alBRyGHxuDYU').hide()
            } else {
                $('#alBRyGHxuDYU').show()
            }
        }
        myFnL()
        // 移动端
        function myH5() {
            let ret = sessionStorage.getItem('isShowH5')
            if (ret == 'hide') {
                $('.advertise').hide()
            } else {
                $('.advertise').show()
            }
        }
  		myH5()
});
