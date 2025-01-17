import {
  DataFrame,
  DataFrameDTO,
  DataTransformContext,
  Field,
  FieldType,
  toDataFrame,
  toDataFrameDTO,
} from '@grafana/data';

import { ModelType, RegressionTransformer, RegressionTransformerOptions } from './regression';

describe('Regression transformation', () => {
  it('it should predict a linear regression to exactly fit the data when the data is f(x) = x', () => {
    const source = [
      toDataFrame({
        name: 'data',
        refId: 'A',
        fields: [
          { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5] },
          { name: 'value', type: FieldType.number, values: [0, 1, 2, 3, 4, 5] },
        ],
      }),
    ];

    const config: RegressionTransformerOptions = {
      modelType: ModelType.linear,
      predictionCount: 6,
      xFieldName: 'time',
      yFieldName: 'value',
    };

    expect(toEquableDataFrames(RegressionTransformer.transformer(config, {} as DataTransformContext)(source))).toEqual(
      toEquableDataFrames([
        toEquableDataFrame({
          name: 'data',
          refId: 'A',
          fields: [
            { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5], config: {} },
            { name: 'value', type: FieldType.number, values: [0, 1, 2, 3, 4, 5], config: {} },
          ],
          length: 6,
        }),
        toEquableDataFrame({
          name: 'linear regression',
          fields: [
            { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5], config: {} },
            { name: 'value predicted', type: FieldType.number, values: [0, 1, 2, 3, 4, 5], config: {} },
          ],
          length: 6,
        }),
      ])
    );
  });
  it('it should predict a linear regression where f(x) = 1', () => {
    const source = [
      toDataFrame({
        name: 'data',
        refId: 'A',
        fields: [
          { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5] },
          { name: 'value', type: FieldType.number, values: [0, 1, 2, 2, 1, 0] },
        ],
      }),
    ];

    const config: RegressionTransformerOptions = {
      modelType: ModelType.linear,
      predictionCount: 6,
      xFieldName: 'time',
      yFieldName: 'value',
    };

    expect(toEquableDataFrames(RegressionTransformer.transformer(config, {} as DataTransformContext)(source))).toEqual(
      toEquableDataFrames([
        toEquableDataFrame({
          name: 'data',
          refId: 'A',
          fields: [
            { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5], config: {} },
            { name: 'value', type: FieldType.number, values: [0, 1, 2, 2, 1, 0], config: {} },
          ],
          length: 6,
        }),
        toEquableDataFrame({
          name: 'linear regression',
          fields: [
            { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5], config: {} },
            { name: 'value predicted', type: FieldType.number, values: [1, 1, 1, 1, 1, 1], config: {} },
          ],
          length: 6,
        }),
      ])
    );
  });

  it('it should predict a quadratic function', () => {
    const source = [
      toDataFrame({
        name: 'data',
        refId: 'A',
        fields: [
          { name: 'time', type: FieldType.time, values: [0, 1, 2, 3, 4, 5] },
          { name: 'value', type: FieldType.number, values: [0, 1, 2, 2, 1, 0] },
        ],
      }),
    ];

    const config: RegressionTransformerOptions = {
      modelType: ModelType.polynomial,
      degree: 2,
      predictionCount: 6,
      xFieldName: 'time',
      yFieldName: 'value',
    };

    const result = RegressionTransformer.transformer(config, {} as DataTransformContext)(source);

    expect(result[1].fields[1].values[0]).toBeCloseTo(-0.1, 1);
    expect(result[1].fields[1].values[1]).toBeCloseTo(1.2, 1);
    expect(result[1].fields[1].values[2]).toBeCloseTo(1.86, 1);
    expect(result[1].fields[1].values[3]).toBeCloseTo(1.86, 1);
    expect(result[1].fields[1].values[4]).toBeCloseTo(1.2, 1);
    expect(result[1].fields[1].values[5]).toBeCloseTo(-0.1, 1);
  });
});

function toEquableDataFrame(source: DataFrame): DataFrame {
  return toDataFrame({
    ...source,
    fields: source.fields.map((field: Field) => {
      return {
        ...field,
        config: {},
      };
    }),
  });
}

function toEquableDataFrames(data: DataFrame[]): DataFrameDTO[] {
  return data.map((frame) => toDataFrameDTO(frame));
}
