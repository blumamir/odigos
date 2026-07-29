import { gql } from '@apollo/client';

export const UPDATE_GO_ENTERPRISE_OFFSETS = gql`
  mutation UpdateGoEnterpriseOffsets($content: String!) {
    updateGoEnterpriseOffsets(content: $content)
  }
`;
