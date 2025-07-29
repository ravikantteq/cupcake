// MongoDB initialization script for Cupcake Kafka Test Framework
// This script sets up the initial database structure and indexes

// Switch to the cupcake database
db = db.getSiblingDB('cupcake');

// Create collections with validation schemas
db.createCollection('flows', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'steps', 'createdAt', 'createdBy'],
      properties: {
        name: {
          bsonType: 'string',
          description: 'Flow name is required and must be a string'
        },
        description: {
          bsonType: 'string',
          description: 'Flow description must be a string'
        },
        version: {
          bsonType: 'string',
          description: 'Version must be a string'
        },
        steps: {
          bsonType: 'array',
          items: {
            bsonType: 'object',
            required: ['stepId', 'type', 'config'],
            properties: {
              stepId: { bsonType: 'string' },
              type: { 
                bsonType: 'string',
                enum: ['produce', 'consume', 'validate', 'delay']
              },
              config: { bsonType: 'object' }
            }
          }
        },
        createdAt: {
          bsonType: 'date',
          description: 'Created date is required'
        },
        updatedAt: {
          bsonType: 'date',
          description: 'Updated date must be a date'
        },
        createdBy: {
          bsonType: 'string',
          description: 'Creator is required'
        }
      }
    }
  }
});

db.createCollection('suites', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'flows', 'createdAt'],
      properties: {
        name: {
          bsonType: 'string',
          description: 'Suite name is required'
        },
        description: {
          bsonType: 'string'
        },
        flows: {
          bsonType: 'array',
          items: { bsonType: 'objectId' }
        },
        environment: {
          bsonType: 'string'
        },
        config: {
          bsonType: 'object'
        },
        createdAt: {
          bsonType: 'date'
        },
        updatedAt: {
          bsonType: 'date'
        }
      }
    }
  }
});

db.createCollection('executions', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['suiteId', 'flowId', 'status', 'startTime'],
      properties: {
        suiteId: { bsonType: 'objectId' },
        flowId: { bsonType: 'objectId' },
        status: {
          bsonType: 'string',
          enum: ['running', 'completed', 'failed', 'cancelled']
        },
        startTime: { bsonType: 'date' },
        endTime: { bsonType: 'date' },
        steps: {
          bsonType: 'array',
          items: {
            bsonType: 'object',
            properties: {
              stepId: { bsonType: 'string' },
              status: { bsonType: 'string' },
              input: { bsonType: 'object' },
              output: { bsonType: 'object' },
              errors: {
                bsonType: 'array',
                items: { bsonType: 'string' }
              },
              duration: { bsonType: 'number' }
            }
          }
        },
        metrics: {
          bsonType: 'object',
          properties: {
            totalDuration: { bsonType: 'number' },
            messagesProduced: { bsonType: 'number' },
            messagesConsumed: { bsonType: 'number' },
            errorsCount: { bsonType: 'number' }
          }
        }
      }
    }
  }
});

db.createCollection('consumers', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['groupId', 'topics', 'status', 'createdAt'],
      properties: {
        groupId: {
          bsonType: 'string',
          description: 'Consumer group ID is required'
        },
        topics: {
          bsonType: 'array',
          items: { bsonType: 'string' },
          description: 'Topics array is required'
        },
        status: {
          bsonType: 'string',
          enum: ['active', 'inactive', 'error'],
          description: 'Status must be one of: active, inactive, error'
        },
        config: {
          bsonType: 'object',
          properties: {
            autoOffsetReset: { bsonType: 'string' },
            enableAutoCommit: { bsonType: 'bool' },
            maxPollRecords: { bsonType: 'number' }
          }
        },
        lastHeartbeat: { bsonType: 'date' },
        createdAt: { bsonType: 'date' }
      }
    }
  }
});

db.createCollection('messages', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['topic', 'partition', 'offset', 'timestamp'],
      properties: {
        topic: { bsonType: 'string' },
        partition: { bsonType: 'number' },
        offset: { bsonType: 'number' },
        key: { bsonType: 'string' },
        value: {},
        headers: { bsonType: 'object' },
        timestamp: { bsonType: 'date' },
        consumerGroupId: { bsonType: 'string' },
        executionId: { bsonType: 'objectId' }
      }
    }
  }
});

// Create indexes for better performance
// Flows collection indexes
db.flows.createIndex({ name: 1 }, { unique: true });
db.flows.createIndex({ createdBy: 1 });
db.flows.createIndex({ createdAt: -1 });
db.flows.createIndex({ 'steps.type': 1 });

// Suites collection indexes
db.suites.createIndex({ name: 1 }, { unique: true });
db.suites.createIndex({ environment: 1 });
db.suites.createIndex({ createdAt: -1 });

// Executions collection indexes
db.executions.createIndex({ suiteId: 1, flowId: 1 });
db.executions.createIndex({ status: 1 });
db.executions.createIndex({ startTime: -1 });
db.executions.createIndex({ 'metrics.totalDuration': 1 });

// Consumers collection indexes
db.consumers.createIndex({ groupId: 1 }, { unique: true });
db.consumers.createIndex({ status: 1 });
db.consumers.createIndex({ lastHeartbeat: -1 });

// Messages collection indexes
db.messages.createIndex({ topic: 1, partition: 1, offset: 1 }, { unique: true });
db.messages.createIndex({ timestamp: -1 });
db.messages.createIndex({ consumerGroupId: 1 });
db.messages.createIndex({ executionId: 1 });

// Create TTL index for message cleanup (keep messages for 30 days)
db.messages.createIndex({ timestamp: 1 }, { expireAfterSeconds: 2592000 });

// Create default test data
print('Creating default test data...');

// Insert sample flow
const sampleFlowId = ObjectId();
db.flows.insertOne({
  _id: sampleFlowId,
  name: 'Sample Order Processing Flow',
  description: 'Test flow for order processing system',
  version: '1.0.0',
  steps: [
    {
      stepId: 'step-1',
      type: 'produce',
      config: {
        topic: 'orders-input',
        message: {
          orderId: 'uuid()',
          customerId: 'string()',
          amount: 'number(min=1, max=1000)',
          timestamp: 'timestamp()'
        },
        timeout: 5000
      }
    },
    {
      stepId: 'step-2',
      type: 'consume',
      config: {
        topic: 'orders-processed',
        timeout: 10000,
        expectedCount: 1
      }
    },
    {
      stepId: 'step-3',
      type: 'validate',
      config: {
        expectedMessage: {
          orderId: 'match(step-1.orderId)',
          status: 'enum(processed,completed)',
          processedAt: 'timestamp()',
          totalAmount: 'number()'
        }
      }
    }
  ],
  createdAt: new Date(),
  updatedAt: new Date(),
  createdBy: 'system'
});

// Insert sample suite
db.suites.insertOne({
  name: 'E2E Order Processing Suite',
  description: 'End-to-end testing suite for order processing system',
  flows: [sampleFlowId],
  environment: 'development',
  config: {
    kafkaBroker: 'localhost:9093',
    consumerGroups: ['cupcake-test-group'],
    timeouts: {
      defaultStepTimeout: 10000,
      suiteTimeout: 300000
    }
  },
  createdAt: new Date(),
  updatedAt: new Date()
});

// Insert sample consumer
db.consumers.insertOne({
  groupId: 'cupcake-test-group',
  topics: ['orders-processed', 'payments-processed', 'notifications-sent'],
  status: 'inactive',
  config: {
    autoOffsetReset: 'earliest',
    enableAutoCommit: true,
    maxPollRecords: 100
  },
  createdAt: new Date()
});

print('MongoDB initialization completed successfully!');
print('Collections created: flows, suites, executions, consumers, messages');
print('Indexes created for optimal performance');
print('Sample data inserted for testing');
