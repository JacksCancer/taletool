

class @Scene

	constructor: (@ctx, @log) ->
		@gl = @ctx.gl

	loadShader: (name, file) ->
		$.when(
			$.get "shader/#{file}.vert"
			$.get "shader/#{file}.frag"
			).fail((xhr, status, error) => @log(status + ": " + error))
			.done((vsResult, fsResult) =>
				[vsrc] = vsResult
				[fsrc] = fsResult
				@[name] = @ctx.compile(vsrc, fsrc))

	loadImage: (file) ->
		img = new Image()
		img.src = file
		$(img).imagesLoaded().then(() -> img)

	createTexture: (img) ->
		tex = @gl.createTexture()
		@gl.bindTexture(@gl.TEXTURE_2D, tex)
		#@gl.pixelStorei(@gl.UNPACK_FLIP_Y_WEBGL, true)
		@gl.texImage2D(@gl.TEXTURE_2D, 0, @gl.RGBA, @gl.RGBA, @gl.UNSIGNED_BYTE, img)
		@gl.texParameteri(@gl.TEXTURE_2D, @gl.TEXTURE_MAG_FILTER, @gl.NEAREST)
		@gl.texParameteri(@gl.TEXTURE_2D, @gl.TEXTURE_MIN_FILTER, @gl.LINEAR)
		tex.width = img.width
		tex.height = img.height
		tex

	loadTexture: (name, file) ->
		img = new Image()
		img.src = file
		$(img).imagesLoaded()
			.fail((xhr, status, error) => @log(status + ": " + error))
			.done(() =>
				tex = @gl.createTexture()
				@gl.bindTexture(@gl.TEXTURE_2D, tex)
				#@gl.pixelStorei(@gl.UNPACK_FLIP_Y_WEBGL, true)
				@gl.texImage2D(@gl.TEXTURE_2D, 0, @gl.RGBA, @gl.RGBA, @gl.UNSIGNED_BYTE, img)
				@gl.texParameteri(@gl.TEXTURE_2D, @gl.TEXTURE_MAG_FILTER, @gl.NEAREST)
				@gl.texParameteri(@gl.TEXTURE_2D, @gl.TEXTURE_MIN_FILTER, @gl.LINEAR)
				tex.width = img.width
				tex.height = img.height
				@[name] = tex)

	loadVertexArray: (name, file) ->
		$.getJSON(file, (array) =>
			buf = @gl.createBuffer()
			@gl.bindBuffer(@gl.ARRAY_BUFFER, buf)
			@gl.bufferData(@gl.ARRAY_BUFFER, new Float32Array(array), @gl.STATIC_DRAW)
			buf.length = array.length
			@[name] = buf
			)

	loadData: (name, file) ->
		$.getJSON(file, (data) =>
			@[name] = data
			)

	checkUniform: (shader, name) ->
		@log("unknown uniform " + name) if not shader.uniform[name]?

	checkAttrib: (shader, name) ->
		@log("unknown attribute " + name) if not shader.attrib[name]?

	checkError: () ->
		err = @gl.getError()
		@log("gl error: " + err) if err != @gl.NO_ERROR

