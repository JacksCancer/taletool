



@initGlContext = (canvas, log) ->
	gl = canvas.getContext "experimental-webgl"

	compile = (vsrc, fsrc) ->
		vs = compileShader vsrc, gl.VERTEX_SHADER
		fs = compileShader fsrc, gl.FRAGMENT_SHADER

		if not vs or not fs then return null

		p = gl.createProgram()

		gl.attachShader p, vs
		gl.attachShader p, fs
		gl.linkProgram p

		if not gl.getProgramParameter(p, gl.LINK_STATUS) or not gl.getProgramParameter(p, gl.VALIDATE_STATUS)
			log gl.getProgramInfoLog(p)
			return null

		return {
			program: p
			uniform: bindUniforms(p)
			attrib: bindAttributes(p)
		}


	compileShader = (src, type) ->
		shader = gl.createShader type
		gl.shaderSource shader, src
		gl.compileShader shader

		if not gl.getShaderParameter(shader, gl.COMPILE_STATUS)
			log gl.getShaderInfoLog(shader)
			return null

		return shader

	bindAttributes = (p) ->
		n = gl.getProgramParameter p, gl.ACTIVE_ATTRIBUTES
		attribs = {}
		attribs[gl.getActiveAttrib(p, i).name] = i for i in [0..n-1]
		return attribs

	bindUniforms = (p) ->
		n = gl.getProgramParameter p, gl.ACTIVE_UNIFORMS
		uniforms = {}
		if n > 0
			for i in [0..n-1]
				name = gl.getActiveUniform(p, i).name
				uniforms[name] = gl.getUniformLocation(p, name) 
		return uniforms

	context =
		gl: gl
		compile: compile
		program: null

	context.use = (s) ->
		attribs = {}
		attribs[i] = 1 for a, i of s.attrib
		old = context.program
		if old
			for a, i of old.attrib
				if attribs[i]?
					delete attribs[i]
				else
					attribs[i] = 0			

		for i, mode of attribs
			if mode == 1
				gl.enableVertexAttribArray(i)
			else
				gl.disableVertexAttribArray(i)

		gl.useProgram(s.program)
		context.program = s


	return context


