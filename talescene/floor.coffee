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

		p1 = path[index].v1
		p2 = path[index].v2

		v = vec3.create()
		vec3.sub(v, p2, p1)
		a = vec3.create()
		vec3.sub(a, p0, p1)

		lower = 0
		upper = 1

		for i in [index+1...path.length]
			u1 = @intersection(a, v, p0, path[i].v1)
			u2 = @intersection(a, v, p0, path[i].v2)

			# sort ascending
			if u1 > u2
				[u1, u2] = [u2, u1]

			if u1 >= upper
				return upper

			if u2 <= lower
				return lower

			if u1 > lower
				lower = u1

			if u2 < upper
				upper = u2

		return lower




	calculatePath: (p0, p1, route) ->
		# path is an array of (edge #0, edge #1)
		path = null
		len = NaN
		for i1 in route
			if path?
				e = @adjacentEdge(@triangles[i0], @triangles[i1])
				# v = vec3.create()
				# vec3.lerp(v, @vertices[e[0]], @vertices[e[1]], 0.5)
				# path.push([v, e[0], e[1]])
				path.push({area: i1, v1: @vertices[e[0]], v2: @vertices[e[1]]})
			else
				path = [{pos: vec3.clone(p0), area: i1}]
			i0 = i1

		if path?
			len = 0
			path.push({pos: vec3.clone(p1), area: i1, v1: p1, v2: p1})

			if path.length > 2
				for i in [1...path.length-1]
					s = @optimizeSection(path[i-1].pos, path, i)
					v = vec3.create()
					vec3.lerp(v, path[i].v1, path[i].v2, s)
					path[i].pos = v
					len += vec3.dist(path[i-1].pos, v)

		return { path:path, len:len }


	# calculates an array of vec3s
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