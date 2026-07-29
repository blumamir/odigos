import { gql } from '@apollo/client';

export const GET_GO_ENTERPRISE_OFFSETS = gql`
  query GetGoEnterpriseOffsets {
    goEnterpriseOffsets {
      timestamp
      mods {
        module
        minVersion
        maxVersion
        versions
      }
    }
  }
`;
