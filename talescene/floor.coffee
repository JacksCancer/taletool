class Floor

	constructor: (definition, @log) ->

		# definition is expected to contain a set of vertices
		# and a set of triangles referencing vertices

		@vertices = (vec3.fromValues(def[0], def[1], def[2]) for def in definition.vertices)
		@triangles = (t for t in definition.triangles)
		@connections = new Array(@triangles.length)
		for t0,i0 in @triangles
			@connections[i0] = (i1 for t1,i1 in @triangles when @adjacent(t0, t1))

	adjacent: (t0, t1) ->
		c = 0
		for i0 in t0
			for i1 in t1
				if i0 == i1
					++c
		c == 2

	adjacentEdge: (t0, t1) ->
		(i for i in t0 when i in t1)

	route: (start, end, path) ->

		path.push(start)

		if start == end
			results = [path[..]]
		else
			for i in @connections[start]
				if not (i in path)
					r = @route(i, end, path)
					if r?
						if results?
							r.push.apply(r, results)
						results = r

		path.pop()
		return results

	# optimize route from p0 to p3 constraint to pass between p1 and p2
	# optimizeSection: (result, p0, p1, p2, p3) ->
	# 	v = vec3.create()
	# 	vec3.sub(v, p2, p1)
	# 	d = vec3.create()
	# 	vec3.sub(d, p3, p0)
	# 	vec3.normalize(d, d)
	# 	n = vec3.create()
	# 	vec3.cross(n, d, vec3.cross(n, d, v))
	# 	a = vec3.create()
	# 	vec3.sub(a, p0, p1)
	# 	s = vec3.dot(a, n)
	# 	t = vec3.dot(v, n)

	# 	if s * t < 0
	# 		vec3.copy(result, p1)
	# 		return false
	# 	else if Math.abs(t) < Math.abs(s)
	# 		vec3.add(result, p1, v)
	# 		return false
	# 	else
	# 		vec3.scaleAndAdd(result, p1, v, s/t)
	# 		return true

	# optimize route from p0 to p3 constraint to pass between p1 and p2
	# optimizeSection: (result, p0, p1, p2, p3) ->
	# optimizeSection: (p0, p1, p2, p3) ->
	# 	v = vec3.create()
	# 	vec3.sub(v, p2, p1)
	# 	a = vec3.create()
	# 	vec3.sub(a, p0, p1)

	# 	d = vec3.create()
	# 	vec3.sub(d, p3, p0)
	# 	vec3.normalize(d, d)
	# 	n = vec3.create()
	# 	vec3.cross(n, vec3.cross(n, d, v), d)
		
	# 	s = vec3.dot(a, n)
	# 	t = vec3.dot(v, n)

	# 	# t should be positive since n and v are pointing in the same direction
	# 	throw t if t < 0

	# 	if s <= 0
	# 		return 0
	# 	else if t <= s
	# 		return 1
	# 	else
	# 		return s/t

	# calculates the nearest point q on section p0+a to p0+a+v and section p0 to p1
	# and returns u where q = p0 + a + v * u
	intersection: (a, v, p0, p1) ->
		d = vec3.create()
		n = vec3.create()

		vec3.sub(d, p1, p0)
		dlen = vec3.len(d)

		# a == 0: s = 0
		# d + a == 0
		# d == 0: n = v
		# d x v == 0: n = 0
		# a == v: s = t

		if dlen > 0
			vec3.scale(d, d, 1 / dlen)
			vec3.cross(n, vec3.cross(n, d, v), d)
		else
			vec3.copy(n, v)
		
		s = vec3.dot(a, n)
		t = vec3.dot(v, n)

		if t > 0.000001
			if -t <= s and s <= 2 * t
				if vec3.dot(vec3.scaleAndAdd(n, a, v, -s/t), d) <= 0
					if s < 0
						0
					else if s > t
						1
					else
						s/t
				else if vec3.dot(d, v) > 0
					1
				else
					0
			else if vec3.dot(vec3.cross(n, v, a), vec3.cross(d, d, a)) > 0
				1
			else
				0
		else if vec3.dot(d, v) > 0
			1
		else
			0

	optimizeSection: (p0, path, index) ->

		p1 = path[index][0]
		p2 = path[index][1]

		v = vec3.create()
		vec3.sub(v, p2, p1)
		a = vec3.create()
		vec3.sub(a, p0, p1)

		@log("sections " + index + ", " + vec3.str(p0) + ", " + vec3.str(p1) + ", " + vec3.str(p2))

		bounds = [0, 1]
		u = new Array(2)

		for i in [index+1...path.length]
			for j in [0..1]
				u[j] = @intersection(a, v, p0, path[i][j])

			# sort ascending
			if u[0] > u[1]
				[u[0], u[1]] = [u[1], u[0]]

			@log(i + ": " + bounds[0] + ", " + bounds[1] + " (" + u[0] + ", " + u[1] + ")")

			if u[0] >= bounds[1]
				return bounds[1]

			if u[1] <= bounds[0]
				return bounds[0]

			if u[0] > bounds[0]
				bounds[0] = u[0]

			if u[1] < bounds[1]
				bounds[1] = u[1]

		return bounds[0]




	calculatePath: (p0, p1, route) ->
		# path is an array of triples (position, edge #0, edge #1)
		path = null
		len = NaN
		for i1 in route
			if path?
				e = @adjacentEdge(@triangles[i0], @triangles[i1])
				# v = vec3.create()
				# vec3.lerp(v, @vertices[e[0]], @vertices[e[1]], 0.5)
				# path.push([v, e[0], e[1]])
				path.push([@vertices[e[0]], @vertices[e[1]]])
			else
				path = [p0]
			i0 = i1

		if path?
			len = 0
			path.push([p1, p1])

			if path.length > 2
				for i in [1...path.length-1]
					# a = path[i-1][0]
					# a = 
					# b = path[i][1]
					# b = @vertices[path[i][1]]
					# c = @vertices[path[i][2]]
					# d = path[i+1][0]

					# s = @optimizeSection(a, b, c, d)
					s = @optimizeSection(path[i-1], path, i)
					v = vec3.create()
					vec3.lerp(v, path[i][0], path[i][1], s)
					path[i] = v

					len += vec3.dist(path[i-1], v)

					#j = i + 1
					#@optimizeSection(path[i][0], a, b, c, path[j][0])

					# while j < path.length and 
					# 	j++

			path[path.length-1] = p1

		return { path:path, len:len }


	navigate: (p0, i0, p1, i1) ->
		shortest = null
		minlen = NaN
		for r in @route(i0, i1, [])
			path = @calculatePath(p0, p1, r)
			if not (path.len > minlen)
				minlen = path.len
				shortest = path.path

		return shortest

	pick: (p0, p1) ->
		v = vec3.create()
		vec3.sub(v, p1, p0)

		itriangle = null
		minpick = Infinity
		mindist = Infinity
		nearest = vec3.create()

		a = vec3.create()
		b = vec3.create()
		n = vec3.create()
		x = vec3.create()
		na = vec3.create()
		nb = vec3.create()

		for t, i in @triangles
			p = @vertices[t[0]]
			vec3.sub(a, @vertices[t[1]], p)
			vec3.sub(b, @vertices[t[2]], p)

			vec3.cross(n, a, b)

			q = vec3.dot(v, n)

			if Math.abs(q) < 0.01
				continue

			s = (vec3.dot(p, n) - vec3.dot(p0, n)) / q

			vec3.cross(na, a, n)
			vec3.cross(nb, b, n)

			vec3.scaleAndAdd(x, p0, v, s)
			vec3.sub(x, x, p)

			ra = vec3.dot(x, nb) / vec3.dot(a, nb)
			rb = vec3.dot(x, na) / vec3.dot(b, na)

			r = ra + rb

			if 0 < ra and 0 < rb and r < 1
				if s < minpick
					minpick = s
					itriangle = i
					vec3.add(nearest, p, x)
			else if minpick == Infinity
				# vector from p to edge start
				c = n
				# edge vector
				d = null
				if rb < 0
					d = a
					vec3.set(c, 0, 0, 0)
				else if ra < 0
					d = b
					vec3.set(c, 0, 0, 0)
				else
					d = b
					vec3.sub(d, b, a)
					vec3.sub(x, x, a)
					c = a

				u = vec3.dot(d, x) / vec3.sqrLen(d)
				u = Math.max(0, Math.min(1, u))
				vec3.scale(d, d, u)

				dist = vec3.dist(d, x)
				if dist < mindist
					mindist = dist
					itriangle = i
					vec3.add(d, d, c)
					vec3.add(nearest, p, d)

		return [nearest, itriangle]


		
class FloorRenderer
	constructor: (@graph) ->

	updateGeometry: (gl) ->

		@vertices = new Float32Array(@graph.vertices.length * 3)
		@indices = new Uint16Array(@graph.triangles.length * 3)

		for v, i in @graph.vertices
			for j in [0...3]
				@vertices[i*3 + j] = v[j]

		for t, i in @graph.triangles
			for j in [0...3]
				@indices[i*3 + j] = t[j]
		
		@vertexBuffer = gl.createBuffer()
		@indexBuffer = gl.createBuffer()
		gl.bindBuffer(gl.ARRAY_BUFFER, @vertexBuffer)
		gl.bufferData(gl.ARRAY_BUFFER, @vertices, gl.STATIC_DRAW)
		gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, @indexBuffer)
		gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, @indices, gl.STATIC_DRAW)


	draw: (gl, attrib, i) ->

		gl.bindBuffer(gl.ARRAY_BUFFER, @vertexBuffer)
		gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, @indexBuffer)
		gl.vertexAttribPointer(attrib, 3, gl.FLOAT, false, 0, 0)
		if i?
			gl.drawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, i*@indices.BYTES_PER_ELEMENT*3)
		else
			n = @indices.byteLength / @indices.BYTES_PER_ELEMENT
			gl.drawElements(gl.TRIANGLES, n, gl.UNSIGNED_SHORT, 0)

@Floor = Floor
@FloorRenderer = FloorRenderer