

initScene = (ctx, log) ->

	scene = {}
	gl = ctx.gl

	loadShader = (name, file) ->
		$.when(
			$.get "shader/#{file}.vert"
			$.get "shader/#{file}.frag"
			).fail((xhr, status, error) -> log(status + ": " + error))
			.done((vsResult, fsResult) ->
				#log("loaded "+ vsrc[0] + " and " + fsrc[0])
				[vsrc] = vsResult
				[fsrc] = fsResult
				scene[name] = ctx.compile(vsrc, fsrc))

	loadTexture = (name, file) ->
		img = new Image()
		img.src = file
		$(img).imagesLoaded()
			.fail((xhr, status, error) -> log(status + ": " + error))
			.done(() ->
				#log(img)
				tex = gl.createTexture()
				gl.bindTexture(gl.TEXTURE_2D, tex)
				#gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true)
				gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, img)
				gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
				gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
				tex.width = img.width
				tex.height = img.height
				scene[name] = tex)

	loadVertexArray = (name, file) ->
		$.getJSON(file, (array) ->
			buf = gl.createBuffer()
			gl.bindBuffer(gl.ARRAY_BUFFER, buf)
			gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(array), gl.STATIC_DRAW)
			buf.length = array.length
			scene[name] = buf
			)
	
	scene.vertices = gl.createBuffer()
	gl.bindBuffer(gl.ARRAY_BUFFER, scene.vertices)
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([
		0, 0,
		1, 0,
		0, 1,
		1, 1,
		]), gl.STATIC_DRAW)

	checkUniform = (shader, name) ->
		log("unknown uniform " + name) if not shader.uniform[name]?

	checkAttrib = (shader, name) ->
		log("unknown attribute " + name) if not shader.attrib[name]?		


	scene.initGl = () ->

		width = scene.bg0.width
		height = scene.bg0.height

		scene.fb = gl.createFramebuffer()
		scene.fbTemp = gl.createFramebuffer()
		scene.fbColor = gl.createTexture()
		scene.fbTempColor = gl.createTexture()
		scene.fbTempWidth = width * 8
		scene.fbTempHeight = height * 8
		scene.fbWidth = width * 2
		scene.fbHeight = height * 2

		gl.bindTexture(gl.TEXTURE_2D, scene.fbColor)
		gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, scene.fbWidth, scene.fbHeight, 0, gl.RGBA, gl.UNSIGNED_SHORT_4_4_4_4, null)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		gl.bindTexture(gl.TEXTURE_2D, scene.fbTempColor)
		gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, scene.fbTempWidth, scene.fbTempHeight, 0, gl.RGBA, gl.UNSIGNED_SHORT_4_4_4_4, null)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		gl.bindFramebuffer(gl.FRAMEBUFFER, scene.fb);
		gl.framebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, scene.fbColor, 0)

		if gl.FRAMEBUFFER_COMPLETE != gl.checkFramebufferStatus(gl.FRAMEBUFFER)
			log("fb status: "+ gl.checkFramebufferStatus(gl.FRAMEBUFFER))


		gl.bindFramebuffer(gl.FRAMEBUFFER, scene.fbTemp);
		gl.framebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, scene.fbTempColor, 0)

		if gl.FRAMEBUFFER_COMPLETE != gl.checkFramebufferStatus(gl.FRAMEBUFFER)
			log("fb status: "+ gl.checkFramebufferStatus(gl.FRAMEBUFFER))


		shader = scene.composeshader
		gl.useProgram(shader.program)

		checkUniform(shader, "u_scale")
		checkUniform(shader, "u_translate")
		checkUniform(shader, "u_tex")
		checkUniform(shader, "u_texscale")
		checkAttrib(shader, "a_coord")

		gl.uniform2f(shader.uniform.u_scale, 1, 1)
		gl.uniform2f(shader.uniform.u_translate, 0, 0)
		gl.uniform1i(shader.uniform.u_tex, 0)
		gl.uniform2f(shader.uniform.u_texscale, 1 / scene.fbTempWidth, 1 / scene.fbTempHeight)


		dwscale = scene.fbWidth / gl.drawingBufferWidth
		dhscale = scene.fbHeight / gl.drawingBufferHeight
		dscale = Math.max(dwscale, dhscale)

		shader = scene.dotshader
		gl.useProgram(shader.program)

		checkUniform(shader, "u_tex")
		checkUniform(shader, "u_tex2")
		checkUniform(shader, "u_veca")
		checkUniform(shader, "u_vecb")
		checkUniform(shader, "u_scale")
		checkUniform(shader, "u_texscale")
		checkAttrib(shader, "a_coord")

		gl.uniform2f(shader.uniform.u_veca, 1, 0)
		gl.uniform2f(shader.uniform.u_vecb, Math.sin(Math.PI/6), Math.cos(Math.PI/6))
		gl.uniform2f(shader.uniform.u_scale, dwscale / dscale, dhscale / dscale)
		gl.uniform2f(shader.uniform.u_texscale, 1 / scene.fbWidth, 1 / scene.fbHeight)

		gl.uniform1i(shader.uniform.u_tex, 0)
		gl.uniform1i(shader.uniform.u_tex2, 1)


		wscale = width / scene.fbTempWidth
		hscale = height / scene.fbTempHeight
		scale = Math.max(wscale, hscale)

		shader = scene.blendshader
		gl.useProgram(shader.program)

		checkUniform(shader, "u_scale")
		checkUniform(shader, "u_factor")
		checkUniform(shader, "u_tex0")
		checkUniform(shader, "u_tex1")
		checkAttrib(shader, "a_coord")

		gl.uniform1i(shader.uniform.u_tex0, 0)
		gl.uniform1i(shader.uniform.u_tex1, 1)
		gl.uniform2f(shader.uniform.u_scale, wscale/scale, hscale/scale)

		shader = scene.animshader
		gl.useProgram(shader.program)

		checkUniform(shader, "u_tex0")
		checkUniform(shader, "u_tex1")
		checkUniform(shader, "u_progress")
		checkUniform(shader, "u_transform")
		checkUniform(shader, "u_texscale")
		checkUniform(shader, "u_scale")
		checkAttrib(shader, "a_point")
		checkAttrib(shader, "a_move")
		checkAttrib(shader, "a_size")

		gl.uniform1i(shader.uniform.u_tex0, 0)
		gl.uniform1i(shader.uniform.u_tex1, 1)
		gl.uniform4f(shader.uniform.u_transform, -width, -height, 1 / (scene.fbTempWidth * scale), 1 / (scene.fbTempHeight * scale))
		gl.uniform2f(shader.uniform.u_texscale, 1 / width, 1 / height)
		gl.uniform1f(shader.uniform.u_scale, 1 / scale)

		return


	scene.render = (factor) ->

		gl.bindFramebuffer(gl.FRAMEBUFFER, scene.fbTemp);
		gl.viewport(0, 0, scene.fbTempWidth, scene.fbTempHeight)
		gl.clearColor(0.4, 0.4, 0, 1)
		gl.clear(gl.COLOR_BUFFER_BIT)

		shader = scene.blendshader
		ctx.use(shader)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, scene.bg0)
		gl.activeTexture(gl.TEXTURE1)
		gl.bindTexture(gl.TEXTURE_2D, scene.bg1)

		gl.uniform1f(shader.uniform.u_factor, factor)

		gl.bindBuffer(gl.ARRAY_BUFFER, scene.vertices)
		# gl.enableVertexAttribArray(shader.attrib.a_coord)
		gl.vertexAttribPointer(shader.attrib.a_coord, 2, gl.FLOAT, false, 0, 0)
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4)
		# gl.disableVertexAttribArray(shader.attrib.a_coord)


		shader = scene.animshader
		ctx.use(shader)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, scene.tex0)
		gl.activeTexture(gl.TEXTURE1)
		gl.bindTexture(gl.TEXTURE_2D, scene.tex1)

		gl.enable(gl.BLEND)

		gl.uniform1f(shader.uniform.u_progress, factor)

		gl.bindBuffer(gl.ARRAY_BUFFER, scene.animverts)
		# gl.enableVertexAttribArray(shader.attrib.a_point)
		# gl.enableVertexAttribArray(shader.attrib.a_move)
		# gl.enableVertexAttribArray(shader.attrib.a_size)
		gl.vertexAttribPointer(shader.attrib.a_point, 2, gl.FLOAT, false, 4 * 6, 0)
		gl.vertexAttribPointer(shader.attrib.a_move, 2, gl.FLOAT, false, 4 * 6, 4 * 2)
		gl.vertexAttribPointer(shader.attrib.a_size, 2, gl.FLOAT, false, 4 * 6, 4 * 4)
		gl.drawArrays(gl.POINTS, 0, scene.animverts.length / 6)


		shader = scene.composeshader
		ctx.use(shader)

		gl.bindFramebuffer(gl.FRAMEBUFFER, scene.fb)
		gl.viewport(0, 0, scene.fbWidth, scene.fbHeight)
		# gl.clearColor(0.7, 0.9, 0.9, 1)
		gl.clearColor(1, 1, 1, 1)
		gl.clear(gl.COLOR_BUFFER_BIT)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, scene.fbTempColor)

		gl.enable(gl.BLEND)
		gl.blendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

		gl.bindBuffer(gl.ARRAY_BUFFER, scene.vertices)
		gl.vertexAttribPointer(shader.attrib.a_coord, 2, gl.FLOAT, false, 0, 0)
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4)

		gl.disable(gl.BLEND)


		shader = scene.dotshader
		ctx.use(shader)

		gl.bindFramebuffer(gl.FRAMEBUFFER, null)
		gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, scene.fbColor)

		gl.activeTexture(gl.TEXTURE1)
		gl.bindTexture(gl.TEXTURE_2D, scene.fbTempColor)

		gl.bindBuffer(gl.ARRAY_BUFFER, scene.vertices)
		gl.vertexAttribPointer(shader.attrib.a_coord, 2, gl.FLOAT, false, 0, 0)
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4)

		#range = gl.getParameter(gl.ALIASED_POINT_SIZE_RANGE)
		#log(range[0] + ", " + range[1])

	scene.loading = 
		$.when(
			loadShader("blendshader", "blend")
			loadShader("animshader", "anim")
			loadShader("quadpixshader", "quadpix")
			loadShader("dotshader", "dot")
			loadShader("composeshader", "compose")
			loadTexture("bg0", "sample/bg0.png")
			loadTexture("bg1", "sample/bg1.png")
			loadTexture("tex0", "sample/test1.png")
			loadTexture("tex1", "sample/test2.png")
			loadVertexArray("animverts", "sample/anim.txt")
			)

	return scene
