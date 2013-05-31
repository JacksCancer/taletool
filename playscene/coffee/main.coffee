$ ->
	$.ajaxSetup
		cache: false

	log = (text) ->
		$("#log").append text + "\n"

	ctx = initGlContext($('#talestage').get(0), log)

	scene = new PlayScene(ctx, log)

	$.when(scene.loading).done(() ->
		err = ctx.gl.getError()
		log("error: " + err) unless err == ctx.gl.NO_ERROR
		scene.initGl()
		scene.render(0.0)

		$("#progress").slider(
			slide: (ev, ui) ->
				if scene.loading.state() == "resolved"
					scene.render(ui.value)
			)

		$("#talestage").click((e) ->
			offset = $(this).offset()
			x = e.pageX - offset.left
			y = e.pageY - offset.top
			scene.click(x, y)
			scene.render(0.0)
			)
		$("#talestage").mousemove((e) ->
			offset = $(this).offset()
			x = e.pageX - offset.left
			y = e.pageY - offset.top
			scene.mouse(x, y)
			scene.render(0.0)
			)
		)

	
		

		