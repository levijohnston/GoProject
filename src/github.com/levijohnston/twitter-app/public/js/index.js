$(function(){

	$(document).on("click", "#sign-up", function(){
		console.log("clicked");
		$("#sign-up-box").show();
		$( "body" ).css({"opacity": "0.7"});
		$("body").css({"background": "rgba(0, 0, 0, 0.2)"})
	});



});