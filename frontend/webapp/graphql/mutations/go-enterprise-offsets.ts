import { gql } from '@apollo/client';

export const UPDATE_GO_ENTERPRISE_OFFSETS = gql`
  mutation UpdateGoEnterpriseOffsets($content: String!) {
    updateGoEnterpriseOffsets(content: $content)
  }
`;

export const CHECK_GO_ENTERPRISE_OFFSETS_UPDATES = gql`
  mutation CheckGoEnterpriseOffsetsUpdates($content: String!) {
    checkGoEnterpriseOffsetsUpdates(content: $content) {
      hasUpdates
      currentTimestamp
      proposedTimestamp
      mods {
        module
        isNew
        minVersion
        maxVersion
        minorVersions {
          minorVersion
          isNew
          versions {
            version
            isNew
          }
        }
      }
    }
  }
`;
