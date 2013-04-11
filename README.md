godiscuz
========

GO语言版本的discuz(ucenter)开发包

#简介
 是discuz(ucenter)论坛的整合开发包的go语法版本,该库不依赖其他的,已经实现了核心的算法,当前方法还没有全部实现,因为时间原因,我也暂时
只是用到这些,已经有的方法都是通过测试的,欢迎大家来扩展,可以参考现在的那些方法的例子,挺简单的,只是要测试,挺耗时的.
#安装

    go get -u github.com/iyf/godiscuz
#开始使用

   import "github.com/iyf/godiscuz"
   godiscuz.discuz.Register("appid", "你的key", "ucenter所在url,最后不要有/", "uc")
   discuz.BindGetcreditsettings(creSet) //这些都是绑定api的,传入相应的方法即可
   discuz.BindLogin(login)
	 discuz.BindUpdatecredit(Updatecredit)
	 discuz.BindLogout(logout)
   //调用方法只需要导入包然后使用即可
