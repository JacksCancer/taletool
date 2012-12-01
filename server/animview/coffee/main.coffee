
log = (text) ->
	$("#log").append text + "\n"



$ ->
	$.ajaxSetup
		cache: false

	ctx = initGlContext($('#talestage').get(0), log)

	scene = initScene(ctx, log)

	$("#progress").slider(
		slide: (ev, ui) ->
			if scene.loading.state() == "resolved"
				scene.render(ui.value / 100.0)
		)

	$.when(scene.loading).done(() ->
		err = ctx.gl.getError()
		log("error: " + err) unless err == ctx.gl.NO_ERROR
		scene.initGl()
		scene.render(0.0))

	
		

		