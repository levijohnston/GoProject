$(function(){

	$(document).on("click", "#sign-up", function(){
		$("#sign-up-box").show();
		$( "body" ).css({"opacity": "0.7"});
		$("body").css({"background": "rgba(0, 0, 0, 0.2)"})
	});

  $(document).on("click", '.close', function(){
    $("#sign-up-box").hide();
  });
});