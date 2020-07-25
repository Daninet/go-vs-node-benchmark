const fastify = require('fastify')({ logger: false });

fastify.register(require('fastify-jwt'), {
  secret: 'secret_jwt_key',
  sign: { algorithm: 'HS512' },
  verify: { algorithms: 'HS512' },
});

fastify.register(require('fastify-postgres'), {
  connectionString: `postgres://postgres:postgres@postgres/demo`,
});

fastify.decorate('authenticate', async function(request, reply) {
  try {
    await request.jwtVerify();
  } catch (err) {
    reply.send(err);
  }
});

fastify.get('/json', async (request, reply) => {
  return { hello: 'world' };
});

fastify.get('/login', async (request, reply) => {
  const payload = { id: 101, first_name: 'John', last_name: 'Doe' };
  const token = fastify.jwt.sign({ user: payload });
  reply.send({ token });
});

fastify.get('/jwt', { preValidation: [fastify.authenticate] }, async (request, reply) => {
  return { payload: request.user };
});

async function getFlights (limit) {
  return await fastify.pg.query(
    `SELECT * FROM bookings.flights WHERE aircraft_code = $1 ORDER BY flight_id LIMIT ${limit}`,
    ['321']
  );
}

fastify.get('/flight', { preValidation: [fastify.authenticate] }, async (request, reply) => {
  const { rows } = await getFlights(1);

  return rows;
});

fastify.get('/flights', { preValidation: [fastify.authenticate] }, async (request, reply) => {
  const { rows } = await getFlights(500);

  return rows;
});

fastify.get('/flights_sorted', { preValidation: [fastify.authenticate] }, async (request, reply) => {
  const { rows } = await getFlights(500);

  rows.sort((a, b) => new Date(b.scheduled_departure).getTime() - new Date(a.scheduled_departure).getTime());

  return rows;
});

const start = async () => {
  try {
    await fastify.listen(3001, '0.0.0.0');
    fastify.log.info(`server listening on ${fastify.server.address().port}`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
}

start();
