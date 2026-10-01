import { gql } from '@apollo/client';

export const GET_UNMATCHED_URL_PATHS = gql`
  query GetUnmatchedUrlPaths($namespace: String!, $kind: String!, $name: String!) {
    unmatchedUrlPaths(namespace: $namespace, kind: $kind, name: $name) {
      server {
        path
        count
        containerName
      }
      client {
        path
        count
        containerName
      }
      serverRecommendedRules {
        template
        reason
        segments {
          templateName
          certainty
          examples
        }
      }
      clientRecommendedRules {
        template
        reason
        segments {
          templateName
          certainty
          examples
        }
      }
      existingConfigs {
        containerName
        templates
        default {
          disabled
          skipPolicy {
            skipForNonSuccessCodes
            skipHttpStatusCodes
          }
        }
      }
    }
  }
`;
